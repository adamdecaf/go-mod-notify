// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adamdecaf/gomodnotify/pkg/database"
	"github.com/adamdecaf/gomodnotify/pkg/modfetch"
	"github.com/adamdecaf/gomodnotify/pkg/modparse"
	"github.com/adamdecaf/gomodnotify/pkg/mods"
	"github.com/adamdecaf/gomodnotify/pkg/nonce"
)

// worker:
// - creates own nonce
// - one transaction
//   - grab N projects from `projects` table (by nonce) without 'started_at is not null and finished_at is null'
//   - insert new rows into scrapes
//   - if tx aborts, retry in <1s (w/ new nonce), but only retry once
// - run each scrape concurrently, perhaps with a splayed start (spread over like 45s)
// - update scrapes
// - metric for current scrapes

var (
	defaultWorkerProjectLimit = 1
)

func spawnWorker(repo database.WorkerRepository) {
	t := time.NewTicker(*flagWorkerInterval)
	for {
		select {
		case <-t.C:
			score := nonce.New()
			projects, err := repo.ScrapeableProjects(score, defaultWorkerProjectLimit)
			if err != nil {
				log.Printf("worker: nonce: %d, but ran into problem: %v", score, err)
			}

			for i := range projects {
				// TODO(adam): log or something here?

				f, err := modfetch.New(projects[i].ImportPath, nil) // TODO(adam): BasicAuth goes here
				if err != nil {
					log.Printf("worker: %s modfetch failed: %v", projects[i].ImportPath, err)
					continue
				}
				dir, err := f.Load(mods.Filenames())
				if err != nil {
					log.Printf("worker: %s module load failed: %v", projects[i].ImportPath, err)
					continue
				}
				mods, err := modparse.ParseFiles(dir, mods.Filenames())
				if err != nil {
					log.Printf("worker: %s modparse failed: %v", projects[i].ImportPath, err)
				}

				fmt.Printf("%#v\n", mods)
			}
		}
	}
}
