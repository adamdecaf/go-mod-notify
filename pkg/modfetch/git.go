// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modfetch

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type GitFetcher struct {
	// modname looks like: github.com/FiloSottile/mkcert
	modname string

	auth *BasicAuth
}

// Load returns a tempdir where dependency files were retrieved.
func (f *GitFetcher) Load(_ []string) (string, error) {
	dir, err := ioutil.TempDir("", "gomodnotify")
	if err != nil {
		return "", fmt.Errorf("unable to create temp dir: %v", err)
	}

	options := &git.CloneOptions{
		URL:   fmt.Sprintf("https://%s.git", f.modname),
		Depth: 1,
	}
	if f.auth != nil {
		// TODO(adam): what do here?
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		options.Auth = &http.BasicAuth{
			Username: "abc123", // can be anything (for personal access tokens)
			Password: f.auth.Password,
		}
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Minute)

	if _, err = git.PlainCloneContext(ctx, dir, false, options); err != nil {
		return "", fmt.Errorf("problem cloning %s: %v", f.modname, err)
	}

	return dir, nil
}
