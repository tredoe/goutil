// Copyright 2011 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package starter

import (
	"os"
	"os/signal"
)

func init() {
	signal.Notify(interrupt, os.Interrupt) // CTRL-C
	signal.Notify(kill, os.Kill)
}
