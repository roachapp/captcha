// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"github.com/roachapp/captcha/pkg/util"
	"testing"
)

func TestNew(t *testing.T) {
	c := DefaultGenerator().New()
	if c == "" {
		t.Errorf("expected id, got empty string")
	}
}

func TestVerify(t *testing.T) {
	g := DefaultGenerator()
	id := g.New()
	if g.Verify(id, []byte{0, 0}) {
		t.Errorf("verified wrong captcha")
	}
	id = g.New()
	d := g.CacheStore.Get(id, false) // cheating
	if !g.Verify(id, d) {
		t.Errorf("proper captcha not verified")
	}
}

func TestReload(t *testing.T) {
	g := DefaultGenerator()
	id := g.New()
	d1 := g.CacheStore.Get(id, false) // cheating
	g.Reload(id)
	d2 := g.CacheStore.Get(id, false) // cheating again
	if bytes.Equal(d1, d2) {
		t.Errorf("reload didn't work: %v = %v", d1, d2)
	}
}

func TestRandomDigits(t *testing.T) {
	d1 := util.RandomDigits(10)
	for _, v := range d1 {
		if v > 9 {
			t.Errorf("digits not in range 0-9: %v", d1)
		}
	}
	d2 := util.RandomDigits(10)
	if bytes.Equal(d1, d2) {
		t.Errorf("digits seem to be not random")
	}
}
