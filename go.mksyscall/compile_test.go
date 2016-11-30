// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build compile

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tredoe/shutil/file"
)

func TestCompile(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Fatal("not Windows system")
	}

	syscallDir := os.ExpandEnv("$GOROOT/src/pkg/syscall")
	filesToRename := []string{"zsyscall_windows_386.go", "zsyscall_windows_amd64.go"}
	filesToGenerate := []string{"syscall_windows.go", "security_windows.go"}

	err := exec.Command("go", "build", ".").Run()
	if err != nil {
		t.Fatal(err)
	}

	// Rename Windows calls generated by "mksyscall_windows.pl"

	for _, v := range filesToRename {
		err = os.Rename(filepath.Join(syscallDir, v), filepath.Join(syscallDir, "_"+v))
		if err != nil {
			t.Fatal(err)
		}
	}

	cmd := exec.Command(
		"go.mksyscall",
		"-w",
		"-conv",
		"false",
		filepath.Join(syscallDir, filesToGenerate[0]),
		filepath.Join(syscallDir, filesToGenerate[1]),
	)
	if err = cmd.Run(); err != nil {
		goto _exit
	}
	os.Remove("go.mksyscall.exe")

	// * * *

	// http://code.google.com/p/go/issues/detail?id=4005
	// The name of a module generated by the Perl version has been changed to
	// avoid introducing new API in Go 1.0.x point releases.
	// So I have to rename that function to title case to can be found in the
	// Windows module.
	err = file.ReplaceAtLineN(
		filepath.Join(syscallDir, "z-syscall_windows.go"),
		[]file.ReplacerAtLine{{"\tproc", `"getCurrentProcessId"`, `"GetCurrentProcessId"`}},
		1,
	)
	if err != nil {
		goto _exit
	}

	// * * *

	// Compile Go using the new sytem calls

	if err = os.Chdir(os.ExpandEnv("$GOROOT/src")); err != nil {
		goto _exit
	}

	cmd = exec.Command("all.bat")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	// == State original
_exit:
	if err != nil {
		t.Error(err)
	}

	os.Chdir(syscallDir)

	if err = os.Remove("z-syscall_windows.go+1~"); err != nil {
		t.Error(err)
	}
	for _, v := range filesToGenerate {
		if err = os.Remove("z" + v); err != nil {
			t.Error(err)
		}
	}
	for _, v := range filesToRename {
		if err = os.Rename("_"+v, v); err != nil {
			t.Error(err)
		}
	}
}
