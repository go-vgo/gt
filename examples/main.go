// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
	"fmt"

	"github.com/go-vgo/gt/file"
	"github.com/go-vgo/gt/hset"
)

func main() {
	sha, err := file.Sha("../file/file.go", "sha256")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sha)

	fileSize, err := file.Size("../file/file.go")
	fmt.Println(fileSize, err)

	f, err := file.ReadIo("../file/file.go")
	fmt.Println(f, err)

	s := hset.New()
	s.Add(1)
	fmt.Println(s)
}
