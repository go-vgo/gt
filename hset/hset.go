// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

/*
Package hset implements a hset backed by a hash table.

Structure is thread safe.

References: http://en.wikipedia.org/wiki/set_%28abstract_data_type%29
*/
package hset

import (
	"fmt"
	"strings"
)

// Hset holds elements in go's native map
var itemExists = struct{}{}

// Contains check if items (one or more) are present in the hset.
// All items have to be present in the hset for the method to return true.
// Returns true if no arguments are passed at all,
// i.e. hset is always superhset of empty hset.
func (hset *Hset) Contains(items ...interface{}) bool {
	for _, item := range items {
		if _, contains := hset.items[item]; !contains {
			return false
		}
	}
	return true
}

// Empty returns true if hset does not contain any elements.
func (hset *Hset) Empty() bool {
	// return hset.Size() == 0
	return hset.Len() == 0
}

// String returns a string representation of container
func (hset *Hset) String() string {
	str := "HasHset\n"
	items := []string{}
	for k := range hset.items {
		items = append(items, fmt.Sprintf("%v", k))
	}

	str += strings.Join(items, ", ")
	return str
}
