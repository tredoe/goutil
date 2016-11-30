// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tredoe/goutil/flagplus"
)

// Flags parsed from command line.
var (
	IsLowerCase = flag.Bool("lowercase", false, "to lower case")
	IsUppercase = flag.Bool("uppercase", false, "to upper case")

	Verbose = flag.Bool("v", false, "mode verbose")
	Str     = flag.String("str", "str", "flag String")
)

func main() {
	cmdHello := new(flagplus.Subcommand)
	cmdHello.UsageLine = "hello [-uppercase] NAME"
	cmdHello.Short = "say hello"
	cmdHello.Long = `"hello" prints out hello to given name.`

	cmdHello.Run = func(cmd *flagplus.Subcommand, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Missing required argument: NAME\n")
			os.Exit(2)
		}

		str := "hello " + args[0]
		if *IsUppercase {
			str = strings.ToUpper(str)
		}
		fmt.Println(str)

		if *Verbose {
			fmt.Println("mode verbose")
		}
	}
	cmdHello.AddFlags("uppercase")

	// * * *

	cmdBye := &flagplus.Subcommand{
		UsageLine: "bye [-lowercase] NAME",
		Short:     "say bye",
		Long:      `"bye" prints out bye to given name`,

		Run: func(cmd *flagplus.Subcommand, args []string) {
			if len(args) != 1 {
				fmt.Fprintf(os.Stderr, "Missing required argument: NAME\n")
				os.Exit(2)
			}

			str := "bye " + args[0]
			if *IsLowerCase {
				str = strings.ToLower(str)
			}
			fmt.Println(str)

			if *Verbose {
				fmt.Println("mode verbose")
			}
		},
	}
	cmdBye.AddFlags("lowercase")

	// * * *

	cmd := flagplus.NewCommand("Test the use of sub-command.", cmdHello, cmdBye)
	cmd.AddGlobalFlags("v", "str")
	cmd.Parse()
}
