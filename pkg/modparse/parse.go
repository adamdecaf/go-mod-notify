// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Version struct {
	Major string
	Minor string
	Patch string
}

// TODO(adam): Equal? Ignore -2018... in Patch string

// TODO(adam): sorting ?

func (v *Version) String() string {
	return fmt.Sprintf("v%s.%s.%s", v.Major, v.Minor, v.Patch)
}

type Modules struct {
	// TODO(adam): mu sync.Mutex ?
	versions map[string]*Version
}

// ParseFile ...
func ParseFile(path string) (*Modules, error) {
	if strings.Contains(path, "../") {
		return nil, errors.New("invalid path")
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ParseFile: problem reading %s", path)
	}
	return Parse(bs)
}

// Parse ...
func Parse(data []byte) (*Modules, error) {
	if len(data) == 0 {
		return nil, errors.New("no go.sum data provided")
	}

	line, lineno := make([]byte, 0), 0
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
	return mods, nil
}

// Example: v0.0.0-20180609054337-500bd5b9081b
func semver(v string) (*Version, error) {
	parts := strings.Split(v, ".")
	if len(parts) <= 2 {
		return nil, fmt.Errorf("Unknown semver: %s", v)
	}

	ver := &Version{}
	if strings.HasPrefix(parts[0], "v") {
		ver.Major = strings.TrimPrefix(parts[0], "v")
	}
	ver.Minor = parts[1]
	ver.Patch = strings.TrimSuffix(strings.Join(parts[2:], "."), "/go.mod")
	return ver, nil
}
