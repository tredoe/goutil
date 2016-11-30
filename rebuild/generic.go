// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build plan9

package rebuild

import "log"

// pkgWatcher represents a watcher for the compiled package files.
type pkgWatcher struct {
	Log *log.Logger

	cmdTocompile string
}

// sysWatcher starts the watcher.
func sysWatcher(cmdTocompile string, _ []string, logg *log.Logger) (*pkgWatcher, error) {
	return &pkgWatcher{logg, cmdTocompile}, nil
}
