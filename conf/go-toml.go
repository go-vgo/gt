// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build !toml
// +build !toml

package conf

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

// Init toml file config
func Init(filePath string, config interface{}) error {
	confLock.Lock()
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Toml init os.ReadFile error: ", err)
		return err
	}

	toml.Unmarshal(fileBytes, config)
	confLock.Unlock()

	return nil
}
