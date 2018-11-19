// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseDep(t *testing.T) {
	bs, err := ioutil.ReadFile(filepath.Join("testdata", "Gopkg.lock"))
	if err != nil {
		t.Fatal(err)
	}

	mods := parseDep(bs)
	if !mods.modulesFound() {
		t.Fatal("no modules found")
	}

	if v := len(mods.versions); v != 1 {
		t.Errorf("got %d modules (expected %d)", v, 1)
	}

	if ver, exists := mods.versions["github.com/gonuts/binary"]; exists {
		if v := ver.String(); v != "v0.1.0" {
			t.Errorf("got %s", v)
		}
	} else {
		t.Errorf("module not found")
	}
}
