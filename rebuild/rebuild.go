// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package rebuild enables to monitor compiled packages for then re-compile a
// related command when any package is updated.
//
// This is very useful during the development of a service.
package rebuild

import (
	"errors"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Prefix for log.
const LOG_PREFIX = "[rebuild] "

// Extension of compiled package files.
const _EXT_COMPILED = ".a"

// To log messages.
var Verbose bool

var errWatcher = errors.New("FAIL! could not monitoring")

// Start watches changes in directories of compiled packages. If there
// is any change, enters in the directory of the command to compile, and it is
// sent a signal Restart.
//
// The argument cmdPath is the import path of source code for the command to run
// by gostart.
// Whether the logger is nil then it is created one by default. The logger
// is always returned.
// Each path in pkgTowatch has to be an import or filesystem path found in
// $GOROOT or $GOPATH.
func Start(cmdPath string, l *log.Logger, pkgTowatch ...string) (*pkgWatcher, error) {
	srv := new(pkgWatcher)
	if l == nil {
		srv.Log = log.New(os.Stdout, LOG_PREFIX, log.LstdFlags)
	}

	// Command
	if cmdPath == "" || build.IsLocalImport(cmdPath) {
		return srv, errors.New("FAIL! import path of command can not be local")
	} else {
		pkg, err := build.Import(cmdPath, build.Default.GOPATH, 0)
		if err != nil {
			return srv, fmt.Errorf("FAIL! at getting command directory: %s", err)
		}

		if !pkg.IsCommand() {
			return srv, fmt.Errorf("FAIL! no command: %s", cmdPath)
		}
		cmdPath = pkg.Dir
	}

	// Packages
	pkgTowatch, pkgFiles, err := checkPkgPath(pkgTowatch, srv.Log)
	if err != nil {
		return srv, err
	}

	w, err := sysWatcher(cmdPath, pkgTowatch, srv.Log)
	if err != nil {
		return srv, err
	}

	if USE_KERNEL {
		go w.watcher(pkgTowatch)
	} else {
		if err = w.watcher(pkgFiles); err != nil {
			return srv, err
		}
	}

	if Verbose {
		srv.Log.Print("Start " + _WATCHER_NAME + " watcher for compiled packages")

		for _, p := range pkgTowatch {
			srv.Log.Printf("Watching %q", p)
		}
	}
	return w, nil
}

// compile compiles the service.
func (w *pkgWatcher) compile() error {
	cmd := exec.Command("go", "install")

	cmd.Dir = w.cmdTocompile
	// For debugging
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FAIL! compile: %s: %s", cmd.Args, err)
	}
	return nil
}

// == Utility
//

// checkPkgPath checks the directories of compiled packages taking an import or
// filesystem path.
// Returns "pkgTowatch", the absolute path of compiled packages; and "pkgFiles",
// for systems without kernel subsystem like inotify.
func checkPkgPath(path []string, logg *log.Logger) (pkgTowatch, pkgFiles []string, err error) {
	var hasError bool
	pkgTowatch = make([]string, len(path))

	for i, p := range path {
		pkg, err := build.Import(p, build.Default.GOPATH, build.AllowBinary)
		/*if err != nil {
			logg.Print("FAIL! checkPkgPath: at getting directory: ", err)
			hasError = true
			continue
		}*/
		if pkg.IsCommand() {
			logg.Print("FAIL! no package: ", p)
			hasError = true
			continue
		}

		pkgPath := filepath.Join(pkg.PkgRoot, runtime.GOOS+"_"+runtime.GOARCH, pkg.ImportPath)
		_, err = os.Stat(pkgPath + _EXT_COMPILED)

		if err != nil && os.IsNotExist(err) {
			// Get compiled packages
			switch files, err := filepath.Glob(filepath.Join(pkgPath, "*"+_EXT_COMPILED)); {
			case err != nil:
				logg.Print("FAIL! checkPkgPath: at getting compiled packages: ", err)
				hasError = true
			case len(files) == 0:
				logg.Printf("FAIL! checkPkgPath: no compiled packages in directory %q", p)
				hasError = true
			case !USE_KERNEL: // there is to get all files to watching
				pkgFiles = append(pkgFiles, files...)
			}
		} else {
			pkgPath += _EXT_COMPILED
		}

		if !hasError {
			pkgTowatch[i] = pkgPath
		}
	}

	if hasError {
		return nil, nil, errWatcher
	}
	return
}
