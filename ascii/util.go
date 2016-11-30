// Copyright 2014 Jonas mg
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ascii

// GetTable returns all characters from a list of ranges.
func GetTable(table ...[]Range) []byte {
	chars_int := make([]int, 0)

	for _, eachTable := range table {
		for _, hex := range eachTable {
			for char := hex.Lo; char <= hex.Hi; char += hex.Stride {
				chars_int = append(chars_int, int(char))
			}
		}
	}

	// They are copied to a new variables with the exact length to save memory.
	chars := make([]byte, len(chars_int))
	for i, char := range chars_int {
		chars[i] = byte(char)
	}

	return chars
}
