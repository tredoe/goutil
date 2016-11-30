// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flagplus

import (
	"testing"

	"github.com/tredoe/goutil/cmdutil"
)

func TestSubcommand(t *testing.T) {
	cmdsInfo := []cmdutil.CommandInfo{
		{
			Stderr: "Test the use of sub-command.\n",
		},
		{
			Args: "hello Joe",
			Out:  "hello Joe\n",
		},
		{
			Args: "hello -uppercase Joe",
			Out:  "HELLO JOE\n",
		},
		{
			Args: "bye Joe",
			Out:  "bye Joe\n",
		},
		{
			Args: "bye -lowercase Joe",
			Out:  "bye joe\n",
		},

		{
			Args: "hello -v Bill",
			Out:  "hello Bill\nmode verbose\n",
		},
		{
			Args: "bye -v Bill",
			Out:  "bye Bill\nmode verbose\n",
		},
	}

	err := cmdutil.TestCommand("testdata", cmdsInfo)
	if err != nil {
		t.Fatal(err)
	}
}
