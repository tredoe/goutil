// Copyright 2014 Jonas mg
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ascii

// Is tests whether character is in the specified table of ranges.
func Is(ranges []Range, char byte) bool {
	for _, r := range ranges {
		if char > r.Hi {
			continue
		}
		if char < r.Lo {
			return false
		}
		return (char-r.Lo)%r.Stride == 0
	}
	return false
}

// IsUpper reports whether the character is an upper case letter.
func IsUpper(char byte) bool {
	return 'A' <= char && char <= 'Z'
}

// IsLower reports whether the character is a lower case letter.
func IsLower(char byte) bool {
	return 'a' <= char && char <= 'z'
}

// IsTitle reports whether the character is a title case letter.
/*func IsTitle(char byte) bool {
	if char < 0x80 { // quick ASCII check
		return false
	}
	return Is(Title, char)
}*/

// IsLetter reports whether the character is a letter.
func IsLetter(char byte) bool {
	char &^= 'a' - 'A'
	return 'A' <= char && char <= 'Z'
}

// IsDigit reports whether the character is a decimal digit.
func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// IsSpace reports whether the character is a white space char.
func IsSpace(char byte) bool {
	return Is(space, char)
}
