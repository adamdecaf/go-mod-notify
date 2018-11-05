// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package mods

import (
	"testing"
)

func TestMods__Filenames(t *testing.T) {
	if len(Filenames()) == 0 {
		t.Error("no filenames")
	}
}
