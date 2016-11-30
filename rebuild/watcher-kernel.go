// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build !plan9

package rebuild

import "github.com/tredoe/goutil/starter"

const (
	_WATCHER_NAME = "kernel"
	USE_KERNEL    = true
)

// Close closes the watcher.
func (w *pkgWatcher) Close() {
	w.watch.Close()

	if Verbose {
		w.Log.Print("Close watcher")
	}
}

// watcher watches compiled packages in operating systems with a kernel
// subsystem for monitoring file systems events.
func (w *pkgWatcher) watcher(pkgTowatch []string) error {
L:
	for { // to logging possible errors
		select {
		case <-starter.Restart:
			break L // compile
		case ev := <-w.watch.Event: // a package has been compiled
			if Verbose {
				w.Log.Print("event: ", ev)
			}
			break L // compile

		case err := <-w.watch.Error:
			w.Log.Print("FAIL! watcher: ", err)
		}
	}

	if err := w.compile(); err != nil {
		w.Log.Print("FAIL! watcher: ", err)
		starter.Stop <- true
	} else {
		starter.Restart <- true
	}
	return nil
}
