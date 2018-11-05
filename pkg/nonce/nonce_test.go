// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package nonce

import (
	"testing"
)

func TestNonce(t *testing.T) {
	results := make([]uint32, 10)
	for i := range results {
		n := New()
		results[i] = n
		if i >= 1 && results[i-1] == n {
			t.Errorf("results[%d] == results[%d], was %d", i-1, i, n)
		}
	}
}
