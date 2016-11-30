// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package data

// +sys getFoo_1foo int, mode *uint32) (err error)
// +sys getFoo_2(foo int, mode *uint32 (err error)
// +sys getFoo_3(foo int, mode *uint32) err error)
// +sys getFoo_4(foo int, mode *uint32) (err error

// +sys getFoo_5(foo int, mode *uint32) (err error = GetFoo
// +sys getFoo_6(foo int, mode *uint32) (err error) =
// +sys getFoo_7(foo int, mode *uint32) (err error)  GetFoo

// +sys getFoo_8(foo int mode *uint32) (err error)
// +sys getFoo_9(foo) (err error)
// +sys getFoo_10(foo int mode *uint32) (err)
