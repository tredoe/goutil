// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rebuild

import (
	"log"

	"golang.org/x/exp/winfsnotify"
)

// pkgWatcher represents a watcher for the compiled package files.
type pkgWatcher struct {
	watch *winfsnotify.Watcher
	Log   *log.Logger

	cmdTocompile string
}

// sysWatcher starts the watcher.
func sysWatcher(cmdTocompile string, pkgTowatch []string, logg *log.Logger) (*pkgWatcher, error) {
	watcher, err := winfsnotify.NewWatcher()
	if err != nil {
		logg.Print("FAIL! sysWatcher: ", err)
		return nil, errWatcher
	}

	ok := true
	// Watch every path
	for _, path := range pkgTowatch {
		if err = watcher.AddWatch(path, winfsnotify.FS_MODIFY); err != nil {
			logg.Print("FAIL! sysWatcher: ", err)
			ok = false
		}
	}

	if !ok {
		return nil, errWatcher
	}
	return &pkgWatcher{watcher, logg, cmdTocompile}, nil
}
