// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"github.com/pelletier/go-toml"
)

type rawLock struct {
	// From https://github.com/golang/dep/blob/5a7960a4cb82dda7253e8b3f921e394cd3c817d3/lock.go#L50
	Projects []rawLockedProject `toml:"projects"`
}

type rawLockedProject struct {
	Name      string   `toml:"name"`
	Branch    string   `toml:"branch,omitempty"`
	Revision  string   `toml:"revision"`
	Version   string   `toml:"version,omitempty"`
	Source    string   `toml:"source,omitempty"`
	Packages  []string `toml:"packages"`
	PruneOpts string   `toml:"pruneopts"`
	Digest    string   `toml:"digest"`
}

func parseDep(data []byte) *Modules {
	mods := &Modules{versions: make(map[string]*Version, 0)}

	var lockfile rawLock
	if err := toml.Unmarshal(data, &lockfile); err != nil {
		return mods
	}

	for i := range lockfile.Projects {
		proj := lockfile.Projects[i]

		ver, err := semver(proj.Version)
		if err != nil {
			continue
		}
		mods.versions[proj.Name] = ver
	}

	return mods
}
