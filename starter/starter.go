// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package starter enables start, restart, and stop services.
//
// Usage:
//
// 1. The server must use the function ExitStatus before of whatever defer
// statement and Wait at the end of the code. Then, to stop or restart your
// service by some condition, use the channels Stop, Error, or Restart.
//
// 2. The server to run must be called through command starter; so, if it exits
// with the status code defined in RESTART, the command can restart it.
package starter

import (
	"log"
	"os"
	"time"
)

// Prefix for log.
const LOG_PREFIX = "[starter] "

// Exit status codes
const (
	_STOP   = 0
	_ERROR  = 1
	RESTART = 33 // <- r3p3at
)

var (
	// System signals
	interrupt = make(chan os.Signal, 2) // handle CTRL-C
	kill      = make(chan os.Signal, 1)

	// Custom signals
	Error   = make(chan bool, 1)
	Restart = make(chan bool, 1)
	Stop    = make(chan bool, 1)

	exitStatus int = _STOP // by default
)

// ExitStatus exits with the exit status code set at function Wait.
//
// It has to be called using defer and before of other defers so it is the last
// one in be called. This is necessary since os.Exit exits immediately and no
// deferred statments can be run.
func ExitStatus() {
	os.Exit(exitStatus)
}

// Wait waits until to receive any signal of interruption or from channels
// Restart, Error, and Stop to set the corresponding exit status.
//
// At pressing CTRL-C the exit status is set to RESTART; being pressed for
// two times, the exit status is set to _STOP. Those codes are used by the
// command "starter" to know when to restart or fisnish a process.
//
// Whether the logger is nil, then it is set one by default.
func Wait(l *log.Logger, verbose bool) {
	if l == nil {
		l = log.New(os.Stdout, LOG_PREFIX, log.LstdFlags)
		if !verbose {
			l.SetFlags(0)
		}
	}

	select {
	case <-kill:
		l.Print("Interrupted")

	// Handle double CTRL-C
	case <-interrupt:
		select {
		case <-time.After(2 * time.Second):
			exitStatus = RESTART
			l.Print("Re-starting...")
		case <-interrupt:
			l.Print("Interrupted")
		}

	case <-Restart:
		exitStatus = RESTART
		l.Print("Re-starting...")
	case <-Stop:
		l.Print("Stopping...")
	case <-Error:
		exitStatus = _ERROR
		l.Print("Stopping due to error...")
	}
}
