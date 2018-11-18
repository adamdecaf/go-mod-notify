// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adamdecaf/gomodnotify/internal"
	"github.com/adamdecaf/gomodnotify/pkg/database"

	"github.com/moov-io/base/admin"
)

var (
	flagAdminAddr = flag.String("admin.addr", ":9090", "Admin HTTP Bind address")

	flagWorkerCount    = flag.Int("worker.count", 10, "How many workers to spawn")
	flagWorkerInterval = flag.Duration("worker.interval", 60*time.Second, "Time between checks for more scrape requests.")
)

func main() {
	flag.Parse()

	log.Printf("Starting gomodnotify/worker:%s\n", internal.Version)

	// Start admin HTTP server
	adminServer := admin.NewServer(*flagAdminAddr)
	go func() {
		log.Printf("listening on %s", adminServer.BindAddr())
		if err := adminServer.Listen(); err != nil {
			err = fmt.Errorf("problem starting admin http: %v", err)
			log.Fatal(err)
		}
	}()
	defer adminServer.Shutdown()

	// Start worker processes
	if *flagWorkerCount < 0 {
		log.Fatalf("invalid worker count: %d", *flagWorkerCount)
	}

	db, err := database.PostgresFromEnv()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	repo := &database.PostgresRepository{
		DB: db,
	}

	wg := sync.WaitGroup{}
	wg.Add(*flagWorkerCount)
	for i := 0; i < *flagWorkerCount; i++ {
		go spawnWorker(repo)
	}
	wg.Wait() // this never completes

	// TODO(adam): handle C-c signal
}
