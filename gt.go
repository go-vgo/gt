// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package gt

import (
	"log"
)

const (
	version string = "v0.10.0.78, Mount Kailash!"
)

// GetVersion get version
func GetVersion() string {
	return version
}

// Try handler error
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// CheckErr check error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// LogErr println error
func LogErr(err error) {
	if err != nil {
		log.Println("error: ", err)
	}
}
