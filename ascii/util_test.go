// Copyright 2014 Jonas mg
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ascii

import "testing"

func TestGetTable(t *testing.T) {
	allChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.-"

	letter_digit := GetTable(Letter, Digit)
	mark := []byte{'_', '.', '-'}

	chars := make([]byte, len(letter_digit)+len(mark))
	copy(chars, letter_digit)
	copy(chars[len(letter_digit):], mark)

	if allChars != string(chars) {
		t.Error("expected to have the same values")
	}
}
