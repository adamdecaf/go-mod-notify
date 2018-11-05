// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package mods

var (
	depdencyFiles = []string{
		"go.sum",     // Go Modules
		"Gopkg.lock", // dep
	}
)

// Filenames returns all the parsable files for various dependency tools. These files
// are assumed to be the tools' "lock file" - a file with specific versions.
func Filenames() []string {
	filenames := make([]string, len(depdencyFiles))
	copy(filenames, depdencyFiles)
	return filenames
}
