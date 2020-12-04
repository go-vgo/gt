// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
)

// GetName get params name
func GetName(args ...string) (string, string) {
	var (
		cmdName = "/bin/bash"
		// cmdName = os.Getenv("SHELL")
		params = "-c"
	)

	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		params = "/C"
	}

	if len(args) > 0 {
		cmdName = args[0]
	}

	if len(args) > 1 {
		params = args[1]
	}

	return cmdName, params
}

// Run run cmd shell
func Run(str string, args ...string) (string, string, error) {
	cmdName, params := GetName(args...)

	fmt.Println("cmd run: ", cmdName, params, ": ", str)
	cmd := exec.Command(cmdName, params, str)

	var out, e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e

	err := cmd.Run()
	return out.String(), e.String(), err
}

// Exec exex command stdout
func Exec(cmdName string, params ...string) bool {
	cmd := exec.Command(cmdName, params...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("cmd.StdoutPipe error: ", err)
		return false
	}
	err = cmd.Start()
	if err != nil {
		log.Println("cmd.Start error: ", err)
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println("cmd.Wait error: ", err)
	}
	return true
}
