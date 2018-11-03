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
)

type GitFetcher struct {
	// modname looks like: github.com/FiloSottile/mkcert
	modname string
}

// Load returns a tempdir where depdency files are stored.
func (g *GitFetcher) Load() (string, error) {
	dir, err := ioutil.TempDir("", "godepnotify")
	if err != nil {
		return "", fmt.Errorf("unable to create temp dir: %v", err)
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Minute)
	_, err = git.PlainCloneContext(ctx, dir, true, &git.CloneOptions{
		URL:   fmt.Sprintf("https://%s.git", g.modname),
		Depth: 1,
	})
	if err != nil {
		return "", fmt.Errorf("problem cloning %s: %v", g.modname, err)
	}

	return dir, nil
}
