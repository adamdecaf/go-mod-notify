// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modfetch

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var (
	// githubRawFileFormat is the URL format Github's 301 to download a raw file.
	// First param is importPath and second is the filepath.
	//
	// master is the branch.
	githubRawFileFormat = "https://%s/raw/master/%s"
)

type GithubFetcher struct {
	// modname looks like: github.com/FiloSottile/mkcert
	modname string
}

var (
	// Set githubTimeoutUntil to a unix ms time in the future to pause requests.
	// All requests must be done with atomic methods.
	githubTimeoutUntil int64 = 0

	githubHttpClient = &http.Client{ // TODO(adam): need to follow their redirect
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     1 * time.Minute,
		},
		Timeout: 30 * time.Second,
	}
)

// Load returns a tempdir where dependency files were retrieved.
func (f *GithubFetcher) Load(filenames []string) (string, error) {
	// TODO(adam): we could use an inmem bloomfilter (key is f.modname + filenames[i])
	// to avoid extra lookups. Flushed every 24 hours if needed?

	// Check if we're supposed to be paused.
	now := time.Now().UTC()
	if atomic.LoadInt64(&githubTimeoutUntil) > now.Unix() {
		timeout := time.Date(0, time.January, 0, 0, 0, int(githubTimeoutUntil), 0, time.UTC) // TODO(adam): 2038 problem
		return "", fmt.Errorf("github: requests paused for %v", now.Sub(timeout))
	}

	dir, err := ioutil.TempDir("", "gomodnotify")
	if err != nil {
		return "", fmt.Errorf("github: unable to create temp dir: %v", err)
	}
	for i := range filenames {
		if err := f.saveFile(dir, filenames[i]); err != nil {
			return "", fmt.Errorf("github: problem saving file %s: %v", filenames[i], err)
		}
	}
	return dir, nil
}

func (f *GithubFetcher) saveFile(dir, filename string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf(githubRawFileFormat, f.modname, filename), nil)
	if err != nil {
		return fmt.Errorf("github: problem making http.Request: %v", err)
	}
	resp, err := githubHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("github: problem making http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 { // Pause for 400, 429, 5xx
		if resp.StatusCode == http.StatusNotFound {
			return nil // try other filenames
		}

		timeout := time.Now().UTC().Add(3 * time.Minute)
		atomic.StoreInt64(&githubTimeoutUntil, timeout.Unix())
		return fmt.Errorf("github: pausing requests until %v", timeout)
	}

	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("github: problem writing %s: %v", path, err)
	}
	defer file.Close()

	n, err := io.Copy(file, resp.Body)
	if err != nil || n == 0 {
		return fmt.Errorf("github: didn't read %s data (n=%d): %v", filename, n, err)
	}
	return nil
}
