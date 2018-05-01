// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package hset

import (
	"encoding/json"
)

// Set hset interface
type Set interface {
	Add(items ...interface{})
	Remove(items ...interface{})
	Clear()
	Contains(items ...interface{}) bool
	Len() int
	Same(other Set) bool
	Values() []interface{}
	String() string
}

// ToJSON outputs the JSON representation of list's elements.
func (set *Hset) ToJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// FromJSON populates list's elements from the input JSON representation.
func (set *Hset) FromJSON(data []byte) error {
	elements := []interface{}{}
	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Clear()
		set.Add(elements...)
	}

	return err
}
