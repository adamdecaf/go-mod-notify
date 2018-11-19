// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
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

func (m *Modules) modulesFound() bool {
	return len(m.versions) > 0
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

	// Attempt parsing with each dependency tool in preference order.
	// Note: Do not change this order as we quit once parsing finds modules.
	if mods := parseGoSum(cp(data)); mods.modulesFound() {
		return mods, nil
	}
	if mods := parseDep(cp(data)); mods.modulesFound() {
		return mods, nil
	}

	return nil, errors.New("No Go dependency management files found")
}

func cp(data []byte) []byte {
	bs := make([]byte, len(data), len(data))
	copy(bs, data)
	return bs
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
