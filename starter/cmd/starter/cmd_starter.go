// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Command starter start, restart, and stop services built in Go with package
// "github.com/tredoe/goutil/starter".
//
// When a service is started, it is created a file with its process identifier
// in the directory given by the environment variable "STARTIT_DIR_PID" if it is
// set, else in the temporary directory.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/tredoe/goutil/starter"
)

var (
	fRestart = flag.Bool("restart", false, "restart the service")
	fStart   = flag.Bool("start", true, "start the sevice (default)")
	fStatus  = flag.Bool("status", false, "to know whether the service is running")
	fStop    = flag.Bool("stop", false, "stop the service")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("FAIL: ")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Tool to start, restart, and stop services.

Usage: starter [option] <service_name> [service_args]
`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NFlag() > 1 || flag.NArg() == 0 {
		usage()
	}

	args := flag.Args()
	service := filepath.Base(args[0])

	pidDir := os.Getenv("STARTIT_DIR_PID")
	if pidDir == "" {
		pidDir = os.TempDir()
	}
	pidFile := filepath.Join(pidDir, service+".pid")

	if *fRestart || *fStop {
		// Get PID from file
		pfile, err := os.Open(pidFile)
		if err != nil {
			log.Fatal(err)
		}
		pinfo, err := pfile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		pid := make([]byte, pinfo.Size())
		if _, err = pfile.Read(pid); err != nil {
			log.Fatal(err)
		}

		pid_i, err := strconv.Atoi(string(pid))
		if err != nil {
			log.Fatal(err)
		}
		proc, err := os.FindProcess(pid_i)
		if err != nil {
			log.Fatal(err)
		}

		if *fRestart {
			if err = proc.Signal(os.Interrupt); err != nil {
				log.Fatal(err)
			}
			fmt.Printf(" * Restarting %s service\n", service)
		} else if *fStop {
			if err = proc.Signal(syscall.SIGTERM); err != nil {
				log.Fatal(err)
			}
			fmt.Printf(" * Stopping %s service\n", service)
		}

	} else if *fStatus {
		_, err := os.Stat(pidFile)
		if os.IsNotExist(err) {
			fmt.Printf("%s: no running\n", service)
		} else {
			fmt.Printf("%s: running\n", service)
		}

	} else if *fStart {
		if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
			log.Fatalf("%s service is already running", service)
		}

		// Create file
		pfile, err := os.OpenFile(pidFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = pfile.Close(); err != nil {
				log.Print(err)
			}
		}()

		isFirst := true
		fmt.Printf(" * Starting %s service\n", service)

		for {
			cmd, err := startCmd(args[0], args[1:])
			if err != nil {
				pfile.Close()
				log.Fatal(err)
			}

			// Write PID to file
			if !isFirst {
				if _, err = pfile.Seek(0, os.SEEK_SET); err != nil {
					pfile.Close()
					log.Fatal(err)
				}
			}
			num, err := pfile.WriteString(strconv.Itoa(cmd.Process.Pid))
			if err != nil {
				pfile.Close()
				log.Fatal(err)
			}
			if !isFirst {
				if err = pfile.Truncate(int64(num)); err != nil {
					pfile.Close()
					log.Fatal(err)
				}
			}

			// Run command
			err = cmd.Wait()
			if msg, ok := err.(*exec.ExitError); ok { // there is error code
				exitStatus := msg.Sys().(syscall.WaitStatus).ExitStatus()
				if exitStatus != starter.RESTART {
					os.Exit(exitStatus)
				}
			} else {
				if err = os.Remove(pidFile); err != nil {
					log.Print(err)
				}
				os.Exit(0)
			}

			isFirst = false
		}
	}
}

// startCmd starts a program.
func startCmd(name string, arg []string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not execute: %s\n%s", cmd.Args, err)
	}
	return cmd, nil
}
