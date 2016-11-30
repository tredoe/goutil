// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build plan9

package rebuild

import (
	"os"
	"time"

	"github.com/tredoe/goutil/starter"
)

const (
	_WATCHER_NAME = "generic"
	USE_KERNEL    = false
)

// Check the compiled packages every x seconds.
var TimeChecking time.Duration = 15 * time.Second

// Close closes the watcher.
func (w *pkgWatcher) Close() {
	if Verbose {
		w.Log.Print("Close watcher")
	}
}

// watcher watches compiled packages in operating systems without kernel
// subsystem for monitoring files, checking its times from time to time.
func (w *pkgWatcher) watcher(pkgFiles []string) error {
	// First checking of times to know if they can be got.
	firstTimes, err := w.mtime(pkgFiles)
	if err != nil {
		return err
	}

	go func(files []string, oldTimes []time.Time) {
		for {
			time.Sleep(TimeChecking)
			hasError := false

			newTimes, err := w.mtime(files)
			if err != nil {
				hasError = true
			}

			if !hasError {
				for i, newT := range newTimes {
					if !newT.Equal(oldTimes[i]) {
						if Verbose {
							w.Log.Printf("event: %q: Mtime", files[i])
						}

						if err := w.compile(); err != nil {
							w.Log.Print("FAIL! watcher: ", err)
							starter.Stop <- true
							return
						}
						starter.Restart <- true
						//close(Restart)
						return
					}
				}
			}
		}
	}(pkgFiles, firstTimes)

	return nil
}

// mtime returns slice with "mtime" for the given files.
func (w *pkgWatcher) mtime(files []string) ([]time.Time, error) {
	var hasError bool
	mtimes := make([]time.Time, len(files))

	for i, v := range files {
		stat, err := os.Stat(v)
		if err != nil {
			w.Log.Print("FAIL! mtime: ", err)
			hasError = true
		} else {
			mtimes[i] = stat.ModTime()
		}
	}

	if hasError {
		return nil, errWatcher
	}
	return mtimes, nil
}
