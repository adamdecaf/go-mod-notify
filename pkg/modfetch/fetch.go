// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modfetch

import (
	"errors"
	"fmt"
	"strings"
)

type Fetcher interface {
	// Load takes tries to download the dependency information to a temp directory. Only files matching
	// filenames are expected to exist if the depdency itself contains the files.
	//
	// This temp directory could have any supported depdency file and it could have multiple.
	//
	// Callers are expected to cleanup the temp directory when finished.
	Load(filenames []string) (string, error)
}

type BasicAuth struct {
	Username string
	Password string
}

// New returns a Fetcher based on the import path
func New(importPath string, auth *BasicAuth) (Fetcher, error) {
	// TODO(adam): handle relative paths
	parts := strings.Split(importPath, "/")
	if len(parts) <= 1 {
		return nil, fmt.Errorf("unknown import path: %s", importPath)
	}
	switch strings.ToLower(parts[0]) {
	case "github.com":
		if auth == nil {
			// No auth provided, so we assume it's a public repo.
			return &GithubFetcher{importPath}, nil
		}
		return &GitFetcher{importPath, auth}, nil
	default:
		return &EmptyFetcher{}, nil
	}
}

type EmptyFetcher struct{}

func (f *EmptyFetcher) Load(_ []string) (string, error) {
	return "", errors.New("EmptyLoader - nil")
}
