// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/adamdecaf/godepnotify/internal"

	"github.com/moov-io/base/admin"
)

var (
	flagAdminAddr = flag.String("admin.addr", ":9090", "Admin HTTP Bind address")

	flagWorkerCount    = flag.Int("worker.count", 10, "How many workers to spawn")
	flagWorkerInterval = flag.Duration("worker.interval", 30*time.Second, "Time between checks for more scrape requests.")
)

func main() {
	flag.Parse()

	log.Printf("Starting godepnotify/worker:%s\n", internal.Version)

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
		log.Printf("invalid worker count: %d", *flagWorkerCount)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(*flagWorkerCount)
	for i := 0; i < *flagWorkerCount; i++ {
		go spawnWorker()
	}
	wg.Wait() // this never completes

	// TODO(adam): handle C-c signal
}
