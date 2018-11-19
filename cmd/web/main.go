// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoModNotify/go-mod-notify"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/admin"
	moovhttp "github.com/moov-io/base/http"
)

var (
	flagHttpAddr  = flag.String("http.addr", ":8080", "HTTP Bind address")
	flagAdminAddr = flag.String("admin.addr", ":9090", "Admin HTTP Bind address")
)

func main() {
	flag.Parse()

	log.Printf("Starting go-mod-notify/web:%s\n", godepnotify.Version)

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

	// Setup HTTP handler
	router := mux.NewRouter()
	moovhttp.AddCORSHandler(router)
	addScrapeEndpoint(router)

	// Start HTTP server
	serve := &http.Server{
		Addr:    *flagHttpAddr,
		Handler: router,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:       false,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
		},
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			log.Fatalf("shutdown err=%v", err)
		}
	}
	defer shutdownServer()

	// TLS Certificate/Key
	tlsCertificate, tlsPrivateKey := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")
	serveViaTLS := tlsCertificate != "" && tlsPrivateKey != ""

	// Serve HTTP and block
	if serveViaTLS {
		log.Printf("transport=HTTPS addr=%v", *flagHttpAddr)
		if err := serve.ListenAndServeTLS(tlsCertificate, tlsPrivateKey); err != nil {
			log.Fatalf("transport=HTTPS err=%v", err)
		}
	} else {
		log.Printf("transport=HTTP addr=%v", *flagHttpAddr)
		if err := serve.ListenAndServe(); err != nil {
			log.Fatalf("transport=HTTP err=%v", err)
		}
	}
}

// moovhttp.GetRequestId(r)
// moovhttp.GetUserId(r)

// moovhttp.InternalError(w, err)
// moovhttp.Problem(w, err)
