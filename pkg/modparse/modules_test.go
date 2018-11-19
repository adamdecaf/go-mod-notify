// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package modparse

import (
	"testing"
)

func TestParse(t *testing.T) {
	cases := map[string]*Version{
		`github.com/DHowett/go-plist v0.0.0-20180609054337-500bd5b9081b h1:WFNhl1+1ofCWWdNFEhut77cmuMXjJYYvkEVloDdaUCI=`:               {"0", "0", "0-20180609054337-500bd5b9081b"},
		`github.com/DHowett/go-plist v0.0.0-20180609054337-500bd5b9081b/go.mod h1:5paT5ZDrOm8eAJPem2Bd+q3FTi3Gxm/U4tb2tH8YIUQ=`:        {"0", "0", "0-20180609054337-500bd5b9081b"},
		`golang.org/x/net v0.0.0-20180627171509-e514e69ffb8b h1:oXs/nlnyk1ue6g+mFGEHIuIaQIT28IgumdSIRMq2aJY=`:                          {"0", "0", "0-20180627171509-e514e69ffb8b"},
		`golang.org/x/net v0.0.0-20180627171509-e514e69ffb8b/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=`:                   {"0", "0", "0-20180627171509-e514e69ffb8b"},
		`golang.org/x/text v0.3.0 h1:g61tztE5qeGQ89tm6NTjjM9VPIm088od1l6aSorWRWg=`:                                                     {"0", "3", "0"},
		`golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=`:                                              {"0", "3", "0"},
		`software.sslmate.com/src/go-pkcs12 v0.0.0-20180114231543-2291e8f0f237 h1:iAEkCBPbRaflBgZ7o9gjVUuWuvWeV4sytFWg9o+Pj2k=`:        {"0", "0", "0-20180114231543-2291e8f0f237"},
		`software.sslmate.com/src/go-pkcs12 v0.0.0-20180114231543-2291e8f0f237/go.mod h1:/xvNRWUqm0+/ZMiF4EX00vrSCMsE4/NHb+Pt3freEeQ=`: {"0", "0", "0-20180114231543-2291e8f0f237"},
	}
	for k, v := range cases {
		mods, err := Parse([]byte(k))
		if err != nil {
			t.Errorf("INPUT: %q ERROR: %v", k, err)
		}
		if v := len(mods.versions); v != 1 {
			t.Errorf("len(mods.versions)=%d", v)
		}
		for path, ver := range mods.versions {
			if ver.Major != v.Major {
				t.Errorf("%s got %s expected %s", path, ver.Major, v.Major)
			}
			if ver.Minor != v.Minor {
				t.Errorf("%s got %s expected %s", path, ver.Minor, v.Minor)
			}
			if ver.Patch != v.Patch {
				t.Errorf("%s got %s expected %s", path, ver.Patch, v.Patch)
			}
		}
	}
}

func TestParseFile(t *testing.T) {
	mods, err := ParseFiles("testdata", []string{"go.sum"})
	if err != nil {
		t.Fatal(err)
	}
	if v := len(mods.versions); v != 4 {
		t.Errorf("got %d modules, expected %d", v, 4)
	}

	check := func(t *testing.T, path string) {
		t.Helper()
		_, exists := mods.versions[path]
		if !exists {
			t.Errorf("didn't find %s as a dependency", path)
		}
	}

	check(t, "github.com/DHowett/go-plist")
	check(t, "golang.org/x/net")
	check(t, "golang.org/x/text")
	check(t, "software.sslmate.com/src/go-pkcs12")
}
