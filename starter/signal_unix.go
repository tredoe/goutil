// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build !windows

package starter

import (
	"os"
	"os/signal"
	"syscall"
)

func init() {
	signal.Notify(interrupt, os.Interrupt) // CTRL-C
	signal.Notify(kill, syscall.SIGTERM)   // kill
	// Note: "kill -9" sends a signal SIGKILL which can not be caught.
}
