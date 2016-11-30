// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"go/scanner"
	"os"
	"testing"
)

const NUM_ERRORS = 10

func Test(t *testing.T) {
	_, err := processFile("testdata/data.go")
	e, ok := err.(scanner.ErrorList)
	if !ok {
		t.Fatal(err)
	}

	if len(e) != NUM_ERRORS {
		t.Errorf("expected %d errors", NUM_ERRORS)
		scanner.PrintError(os.Stderr, err)
	}
}
