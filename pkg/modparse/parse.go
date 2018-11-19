// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func (m *Modules) ForEach(f func(string, *Version)) {
	for path, ver := range m.versions {
		f(path, ver)
	}
}

// ParseFiles returns the first Modules object parsed from all available dependency tools.
func ParseFiles(dir string, paths []string) (*Modules, error) {
	for i := range paths {
		if strings.Contains(paths[i], "../") {
			return nil, fmt.Errorf("invalid path %s", paths[i])
		}

		bs, err := ioutil.ReadFile(filepath.Join(dir, paths[i]))
		if err != nil {
			return nil, fmt.Errorf("ParseFile: problem reading %s", paths[i])
		}
		if len(bs) > 0 {
			return Parse(bs)
		}
	}
	return nil, errors.New("unable to find dependency files")
}

// Parse ...
func Parse(data []byte) (*Modules, error) {
	if len(data) == 0 {
		return nil, errors.New("no go.sum data provided")
	}

	var line []byte
	lineno := 0

	mods := &Modules{versions: make(map[string]*Version)}
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
		line = bytes.TrimSpace(line)

		// skip empty lines
		if len(line) == 0 {
			continue
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
