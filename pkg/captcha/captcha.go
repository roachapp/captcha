// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package captcha implements generation and verification of image and audio
// CAPTCHAs.
//
// A captcha solution is the sequence of digits 0-9 with the defined length.
// There are two captcha representations: image and audio.
//
// An image representation is a PNG-encoded image with the solution printed on
// it in such a way that makes it hard for computers to solve it using OCR.
//
// An audio representation is a WAVE-encoded (8 kHz unsigned 8-bit) sound with
// the spoken solution (currently in English, Russian, Chinese, and Japanese).
// To make it hard for computers to solve audio captcha, the voice that
// pronounces numbers has random speed and pitch, and there is a randomly
// generated background noise mixed into the sound.
//
// This package doesn't require external files or libraries to generate captcha
// representations; it is self-contained.
//
// To make captchas one-time, the package includes a memory storage that stores
// captcha ids, their solutions, and expiration time. Used captchas are removed
// from the store immediately after calling Verify or VerifyString, while
// unused captchas (user loaded a page with captcha, but didn't submit the
// form) are collected automatically after the predefined expiration time.
// Developers can also provide custom store (for example, which saves captcha
// ids and solutions in database) by implementing Store interface and
// registering the object with SetCustomStore.
//
// Captchas are created by calling New, which returns the captcha id.  Their
// representations, though, are created on-the-fly by calling WriteImage or
// WriteAudio functions. Created representations are not stored anywhere, but
// subsequent calls to these functions with the same id will write the same
// captcha solution. Reload function will create a new different solution for
// the provided captcha, allowing users to "reload" captcha if they can't solve
// the displayed one without reloading the whole page.  Verify and VerifyString
// are used to verify that the given solution is the right one for the given
// captcha id.
//
// Server provides an http.Handler which can serve image and audio
// representations of captchas automatically from the URL. It can also be used
// to reload captchas.  Refer to Server function documentation for details, or
// take a look at the example in "capexample" subdirectory.
package captcha

import (
	"bytes"
	"context"
	"errors"
	"github.com/roachapp/captcha/pkg/store"
	"github.com/roachapp/captcha/pkg/util"
	"io"
	"time"
)

var (
	ErrNotFound = errors.New("captcha: id not found")
)

type Generator struct {
	DigitLen int // default 3
	Width int // default 160
	Height int // default 80
	CacheStore store.Store
	PgStore store.Store
}

// New creates a new captcha with the standard length, saves it in the internal
// storage and returns its id.
func (g *Generator) New() string {
	return g.NewLen(g.DigitLen)
}

// NewLen is just like New, but accepts length of a captcha solution as the
// argument.
func (g *Generator) NewLen(length int) string {
	id := util.RandomId()
	digits := util.RandomDigits(length)
	g.CacheStore.Set(id, digits)
	g.PgStore.Set(id, digits)
	return id
}

// Reload generates and remembers new digits for the given captcha id.  This
// function returns false if there is no captcha with the given id.
//
// After calling this function, the image or audio presented to a user must be
// refreshed to show the new captcha representation (WriteImage and WriteAudio
// will write the new one).
func (g *Generator) Reload(id string) bool {
	var old []byte
	if old = g.CacheStore.Get(id, false); old == nil {
		if old = g.PgStore.Get(id, false); old == nil {
			return false
		}
	}

	digits := util.RandomDigits(len(old))
	g.CacheStore.Set(id, digits)
	g.PgStore.Set(id, digits)
	return true
}

// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func (g *Generator) WriteImage(w io.Writer, id string, width, height int) error {
	var d []byte
	if d = g.CacheStore.Get(id, false); d == nil {
		if d = g.PgStore.Get(id, false); d == nil {
			return ErrNotFound
		}
	}

	_, err := util.NewImage(id, d, width, height).WriteTo(w)
	return err
}

// Verify returns true if the given digits are the ones that were used to
// create the given captcha id.
//
// The function deletes the captcha with the given id from the internal
// storage, so that the same captcha can't be verified anymore.
func (g *Generator) Verify(id string, digits []byte) bool {
	if digits == nil || len(digits) == 0 {
		return false
	}

	var realDigits []byte
	cacheDigits := g.CacheStore.Get(id, true)
	pgDigits := g.PgStore.Get(id, true)

	if realDigits = cacheDigits; cacheDigits == nil {
		if realDigits = pgDigits; pgDigits == nil {
			return false
		}
	}

	return bytes.Equal(digits, realDigits)
}

// VerifyString is like Verify, but accepts a string of digits.  It removes
// spaces and commas from the string, but any other characters, apart from
// digits and listed above, will cause the function to return false.
func (g *Generator) VerifyString(id string, digits string) bool {
	if digits == "" {
		return false
	}
	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}
	return g.Verify(id, ns)
}

// DefaultGenerator is used strictly for testing
func DefaultGenerator() *Generator {
	return &Generator{
		CacheStore: store.NewCacheStore(100, 30 * time.Second),
		PgStore: store.NewPostgresStore(context.Background()),
		DigitLen: 3,
		Width: 160,
		Height: 80,
	}
}