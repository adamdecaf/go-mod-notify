// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modfetch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGithubFetcher(t *testing.T) {
	f := &GithubFetcher{"github.com/moov-io/auth"}
	dir, err := f.Load()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Errorf("no temp dir")
	}
	// We'd better see a go.sum file
	if _, err := os.Stat(filepath.Join(dir, "go.sum")); err != nil {
		t.Errorf("couldn't find go.sum: %v", err)
	}
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
}
