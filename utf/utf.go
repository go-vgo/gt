// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package utf

import (
	"fmt"
	"strconv"
	"strings"
)

// CharCodeAt char code at utf-8
func CharCodeAt(str string, n int) rune {
	i := 0
	for _, r := range str {
		if i == n {
			return r
		}
		i++
	}

	return 0
}

// ToUnicode trans string to unicode
func ToUnicode(str string) (uc []string) {
	for _, r := range str {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		uc = append(uc, textUnQ)
	}

	return
}

// Unicode tarans string to unicode string
func Unicode(str string) (u string) {
	for _, v := range ToUnicode(str) {
		u += v
	}
	return
}

// ToUC trans string to Unicode
func ToUC(str string) (uc []string) {
	for _, v := range ToUnicode(str) {
		st := strings.Replace(v, "\\u", "U", -1)
		if st == "\\\\" {
			st = "\\"
		}
		if st == `\"` {
			st = `"`
		}
		uc = append(uc, st)
	}

	return
}

// UnicodeToUTF8 trans Unicode to utf-8
func UnicodeToUTF8(str string) string {
	res := []string{}
	strArr := strings.Split(str, "\\u")
	var snip string
	for _, v := range strArr {
		add := ""
		if len(v) < 1 {
			continue
		}

		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			add = string(rs[4:])
		}

		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			snip += v
		}

		if temp > 0 {
			snip += fmt.Sprintf("%c", temp)
		}
		snip += add
	}

	res = append(res, snip)
	return strings.Join(res, "")
}
