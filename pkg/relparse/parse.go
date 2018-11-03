// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package relparse

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func Parse(_ []byte) {
	resp, err := http.Get("https://github.com/FiloSottile/mkcert/releases.atom")
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	var releases GithubReleases
	if err := xml.Unmarshal(bs, &releases); err != nil {
		log.Printf("ERROR: %v", err)
	}

	for _, entry := range releases.Entry {
		parts := strings.Split(entry.ID, "/")
		if len(parts) == 0 {
			continue
		}

		updated, err := time.Parse("2006-01-02T15:04:05Z", entry.Updated)
		if err != nil {
			continue
		}

		version := parts[len(parts)-1]

		// title := entry.Titel
		// content := entry.Content

		fmt.Printf("%s released %v\n", version, updated)
	}
}
