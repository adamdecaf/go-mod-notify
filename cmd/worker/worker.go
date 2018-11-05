// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"

	"github.com/adamdecaf/godepnotify/pkg/nonce"
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

func spawnWorker() {
	t := time.NewTicker(*flagWorkerInterval)
	for {
		select {
		case <-t.C:
			score := nonce.New()
			log.Printf("worker: using score %d", score)

			// check for new work, execute scrapes, etc
			// ya know, stuff
		}
	}
}
