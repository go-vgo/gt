// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

//go:build !toml
// +build !toml

package conf

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Init toml file config
func Init(filePath string, config interface{}, embed1 ...bool) (err error) {
	confLock.Lock()
	var fileBytes []byte
	if len(embed1) > 0 {
		fileBytes = []byte(filePath)
	} else {
		fileBytes, err = os.ReadFile(filePath)
	}
	if err != nil {
		log.Println("Toml init os.ReadFile error: ", err)
		return err
	}

	err = toml.Unmarshal(fileBytes, config)
	if err != nil {
		return err
	}
	confLock.Unlock()

	return nil
}
