// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package gt

import (
	"log"
	"time"
)

const (
	// Version get version
	Version string = "v0.20.0.153, Mount Kailash!"
)

// GetVersion get version
func GetVersion() string {
	return Version
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

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
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
