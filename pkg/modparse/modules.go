// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"bytes"
	"strings"
)

func parseGoSum(data []byte) *Modules {
	var line []byte
	lineno := 0

	mods := &Modules{versions: make(map[string]*Version, 0)}
	for {
		// Break if we've gone past some really high limit
		lineno++
		if lineno > 50000 {
			break // TODO(adam): log, or something
		}

		// Read for \n
		i := bytes.IndexByte(data, '\n')
		if i < 0 {
			line, data = data, nil
		} else {
			line, data = data[:i], data[i+1:]
		}

		// split at ' '
		parts := strings.Fields(string(line))
		if len(parts) < 2 {
			continue // empty line
		}

		// first two parts should be a path and semver
		path, version := parts[0], parts[1]
		ver, err := semver(version)
		if err != nil {
			continue // TODO(adam): log?
		}
		// TODO(adam): can we assume path doesn't exist already?
		mods.versions[path] = ver
	}
	return mods
}
