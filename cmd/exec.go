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
)

// ExecSh exec command shell
func ExecSh(str string, args ...string) (string, error) {
	var (
		cmdName = "/bin/bash"
		params  = "-c"
	)

	if len(args) > 0 {
		cmdName = args[0]
	}

	if len(args) > 1 {
		params = args[1]
	}

	cmd := exec.Command(cmdName, params, str)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	return out.String(), err
}

// ExecCmd exex command stdout
func ExecCmd(cmdName string, params []string) bool {
	cmd := exec.Command(cmdName, params...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println("cmd.StdoutPipe error: ", err)
		return false
	}
	cmd.Start()

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
	return true
}
