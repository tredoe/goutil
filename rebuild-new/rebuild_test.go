// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rebuild

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/tredoe/goutil/starter"
)

func init() {
	*Verbose = true
}

func TestRebuild(t *testing.T) {
	cmdTocompile := "github.com/tredoe/goutil/rebuild/testdata/testcmd"
	pkgTowatch := "github.com/tredoe/goutil/rebuild/testdata/testpkg"

	// Install command for testing; the package is automatically installed.
	err := os.Chdir(filepath.Join("testdata", "testcmd"))
	if err != nil {
		t.Fatal(err)
	}
	dirCmd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err = exec.Command("go", "install").Run(); err != nil {
		t.Fatal(err)
	}

	// Get full path at installing
	pkg, err := build.Import(pkgTowatch, build.Default.GOPATH, build.FindOnly)
	if err != nil {
		t.Fatal(err)
	}
	cmd := filepath.Join(pkg.BinDir, filepath.Base(cmdTocompile))

	// Initial time of the command
	info, err := os.Stat(cmd)
	if err != nil {
		t.Fatal(err)
	}
	firstTime := info.ModTime()

	// Initialize watcher
	watcher, err := Start(cmdTocompile, nil, pkgTowatch)
	if err != nil {
		t.Fatal("expected start watcher: ", err)
	}
	defer watcher.Close()

	// == Package compilation

	fmt.Printf("\nCompiling...\n")
	go starter.Wait(nil, *Verbose)

	if err = os.Chdir(pkg.Dir); err != nil {
		t.Fatal(err)
	}
	if err = exec.Command("go", "clean", "-i").Run(); err != nil {
		t.Fatal(err)
	}
	if err = exec.Command("go", "install").Run(); err != nil {
		t.Fatal(err)
	}

	if USE_KERNEL {
		time.Sleep(3 * time.Second)
	} else {
		fmt.Printf("Waiting for checking...\n")
		time.Sleep(25 * time.Second)
	}
	starter.Stop <- true

	// Check time of the command
	info2, err := os.Stat(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if firstTime.Equal(info2.ModTime()) {
		t.Error("the command was not re-compiled")
	}

	// == Remove test data

	// directory "pkg.Dir"
	if err = exec.Command("go", "clean", "-i").Run(); err != nil {
		t.Log(err)
	}

	if err = os.Chdir(dirCmd); err != nil {
		t.Log(err)
	} else if err = exec.Command("go", "clean", "-i").Run(); err != nil {
		t.Log(err)
	}
}
