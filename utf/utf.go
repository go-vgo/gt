// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
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

// ToUC trans string to Unicode
func ToUC(str string) (uc []string) {
	for _, r := range str {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		st := strings.Replace(textUnQ, "\\u", "U", -1)
		uc = append(uc, st)
	}

	return uc
}

// UnicodeToUTF8 trans Unicode to utf-8
func UnicodeToUTF8(str string) string {
	i := 0
	if strings.Index(str, `\u`) > 0 {
		i = 1
	}

	strArr := strings.Split(str, `\u`)
	last := len(strArr) - 1
	if len(strArr[last]) > 4 {
		strArr = append(strArr, strArr[last][4:])
		strArr[last] = strArr[last][:4]
	}

	for ; i <= last; i++ {
		if x, err := strconv.ParseInt(strArr[i], 16, 32); err == nil {
			strArr[i] = fmt.Sprintf("%c", x)
		}
	}
	return strings.Join(strArr, "")
}
