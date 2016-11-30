// Copyright 2014 Jonas mg
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ascii

// The representation of a range of ASCII characters. The range runs from Lo to Hi
// inclusive and has the specified stride.
type Range struct {
	Lo     byte
	Hi     byte
	Stride byte
}

// Categories is the set of ASCII data tables.
/*var Categories = map[string][]Range{
	"Nd":      Nd,
	"Lu":      Lu,
	"Ll":      Ll,
	"control": Cc,
	"letter":  letter,
	"space":   space,
}*/

var (
	Digit   = _Nd    // Digit is the set of ASCII digits.
	Upper   = _Lu    // Upper is the set of ASCII upper case letters.
	Lower   = _Ll    // Lower is the set of ASCII lower case letters.
	Control = _Cc    // Control is the set of ASCII control characters.
	Letter  = letter // Letter is the set of ASCII letters.
	Space   = space  // Space is the set of ASCII space characters.
)

var _Nd = []Range{
	Range{0x30, 0x39, 1},
}

var _Lu = []Range{
	Range{0x41, 0x5a, 1},
}

var _Ll = []Range{
	Range{0x61, 0x7a, 1},
}

var _Cc = []Range{
	Range{0x01, 0x1f, 1},
	Range{0x7f, 0x7f, 1},
}

var letter = []Range{
	Range{0x41, 0x5a, 1},
	Range{0x61, 0x7a, 1},
}

var space = []Range{
	Range{0x09, 0x0d, 1},
	Range{0x20, 0x20, 1},
}
