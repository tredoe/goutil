// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/tredoe/osutil"
)

func TestSignal(t *testing.T) {
	go Wait(nil, true)
	time.Sleep(1 * time.Second)
	Restart <- true

	go Wait(nil, true)
	time.Sleep(1 * time.Second)
	Error <- true

	go Wait(nil, false)
	time.Sleep(1 * time.Second)
	Stop <- true
}

func TestCommand(t *testing.T) {
	CMD_MAIN := "starter"
	CMD_TEST := "tester"

	// Build the command and executable test
	err := exec.Command("go", "build", "./"+filepath.Join("cmd", CMD_MAIN)).Run()
	if err != nil {
		t.Fatal(err)
	}
	err = exec.Command("go", "build", "-o", CMD_TEST, "./testdata").Run()
	if err != nil {
		t.Fatal(err)
	}

	CMD_MAIN = "./" + CMD_MAIN

	go func() {
		err = osutil.Exec(CMD_MAIN, "./"+CMD_TEST)
	}()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	if err = osutil.Exec(CMD_MAIN, "-status", CMD_TEST); err != nil {
		t.Error(err)
	}
	time.Sleep(3 * time.Second)
	fmt.Println()

	if err = osutil.Exec(CMD_MAIN, "-restart", CMD_TEST); err != nil {
		t.Error(err)
	}
	time.Sleep(7 * time.Second)
	fmt.Println()

	if err = osutil.Exec(CMD_MAIN, "-stop", CMD_TEST); err != nil {
		t.Error(err)
	}
	time.Sleep(3 * time.Second)
	fmt.Println()

	if err = osutil.Exec(CMD_MAIN, "-status", CMD_TEST); err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)

	// Remove the commands
	for _, v := range []string{CMD_MAIN, CMD_TEST} {
		if err = os.Remove(v); err != nil {
			t.Log(err)
		}
	}
}
