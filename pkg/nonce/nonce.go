// Copyright 2018 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package nonce

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"sync"
	"time"
)

var (
	mrandSetup sync.Once
)

func New() uint32 {
	bs := make([]byte, 4)
	if _, err := rand.Read(bs); err != nil {
		// fallback to math/rand
		mrandSetup.Do(func() {
			mrand.Seed(time.Now().Unix())
		})
		return mrand.Uint32()
	}
	return binary.BigEndian.Uint32(bs)
}
