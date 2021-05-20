// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"github.com/roachapp/captcha/pkg/util"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	s := NewCacheStore(100, 30 * time.Second)
	id := "captcha id"
	d := util.RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestGetClear(t *testing.T) {
	s := NewCacheStore(100, 30 * time.Second)
	id := "captcha id"
	d := util.RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, true)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != nil {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func TestCollect(t *testing.T) {
	//TODO(dchest): can't test automatic collection when saving, because
	//it's currently launched in a different goroutine.
	s := NewCacheStore(10, -1)
	// create 10 ids
	ids := make([]string, 10)
	d := util.RandomDigits(10)
	for i := range ids {
		ids[i] = util.RandomId()
		s.Set(ids[i], d)
	}
	s.(*cacheStore).collect()
	// Must be already collected
	nc := 0
	for i := range ids {
		d2 := s.Get(ids[i], false)
		if d2 != nil {
			t.Errorf("%d: not collected", i)
			nc++
		}
	}
	if nc > 0 {
		t.Errorf("= not collected %d out of %d captchas", nc, len(ids))
	}
}

func BenchmarkSetCollect(b *testing.B) {
	b.StopTimer()
	d := util.RandomDigits(10)
	s := NewCacheStore(9999, -1)
	ids := make([]string, 1000)
	for i := range ids {
		ids[i] = util.RandomId()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			s.Set(ids[j], d)
		}
		s.(*cacheStore).collect()
	}
}
