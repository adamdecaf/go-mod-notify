// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modfetch

import (
	"errors"
	"fmt"
	"log"
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

	// TODO(adam): prometheus metric for time taken to scrape
	// labels: type=$(*modfetch.GithubFetcher | tr -d '*modfetch.')

	var fetcher Fetcher
	switch strings.ToLower(parts[0]) {
	case "github.com":
		if auth == nil {
			// No auth provided, so we assume it's a public repo.
			fetcher = &GithubFetcher{importPath}
		} else {
			fetcher = &GitFetcher{importPath, auth}
		}
	default:
		fetcher = &EmptyFetcher{}
	}
	hasAuth := auth != nil
	log.Printf("using %T for %s dependency retrieval (auth:%v)", fetcher, importPath, hasAuth)
	return fetcher, nil
}

type EmptyFetcher struct{}

func (f *EmptyFetcher) Load(_ []string) (string, error) {
	return "", errors.New("EmptyLoader - nil")
}
