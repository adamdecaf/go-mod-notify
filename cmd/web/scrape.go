// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/GoModNotify/go-mod-notify/pkg/modfetch"
	"github.com/GoModNotify/go-mod-notify/pkg/modparse"
	"github.com/GoModNotify/go-mod-notify/pkg/mods"

	"github.com/gorilla/mux"
	moovhttp "github.com/moov-io/base/http"
)

func addScrapeEndpoint(r *mux.Router) {
	r.Methods("POST").Path("/scrape").HandlerFunc(scrapeEndpoint)
}

func getImportPath(r *http.Request) string {
	return r.URL.Query().Get("importPath")
}

type module struct {
	ImportPath string `json:"importPath"`
	Version    string `json:"version"`
}

func scrapeEndpoint(w http.ResponseWriter, r *http.Request) {
	importPath, requestId := getImportPath(r), moovhttp.GetRequestId(r)
	logger := func(err error) {
		if err != nil && importPath != "" {
			// If we're returning an error let's log that iff X-Request-Id is set
			log.Printf("requestId=%s problem getting modules for %s: %v", requestId, importPath, err)
		}
	}

	if importPath == "" {
		err := errors.New("missing import path")
		logger(err)
		moovhttp.Problem(w, err)
		return
	}

	if requestId != "" {
		log.Printf("requestId=%s getting modules for %s", requestId, importPath)
	}

	// Grab repo
	f, err := modfetch.New(importPath, nil) // TODO(adam): BasicAuth goes here
	if err != nil {
		err = fmt.Errorf("problem grabbing %s: %v", importPath, err)
		logger(err)
		moovhttp.Problem(w, err)
		return
	}
	dir, err := f.Load(mods.Filenames())
	if err != nil {
		err = fmt.Errorf("problem loading %s: %v", importPath, err)
		logger(err)
		moovhttp.Problem(w, err)
		return
	}

	// Find Modules
	mods, err := modparse.ParseFiles(dir, mods.Filenames())
	if err != nil {
		err = fmt.Errorf("problem parsing %s go.sum: %v", importPath, err)
		moovhttp.Problem(w, err)
		return
	}

	// Render json
	var modules []module
	mods.ForEach(func(path string, ver *modparse.Version) {
		modules = append(modules, module{
			ImportPath: path,
			Version:    ver.String(),
		})
	})
	if err = json.NewEncoder(w).Encode(modules); err != nil {
		logger(err)
		moovhttp.Problem(w, err)
		return
	}

	// TODO(adam): instead, write scrape into db (it'll be used as force, and to queue projects)
}
