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
	"fmt"
	"strings"
	"sync"
)

// Hset holds elements in go's native map
type Hset struct {
	items map[interface{}]struct{}
	// items sync.Map[interface{}]struct{}
	sync.RWMutex
}

// Hset holds elements in go's native map
var itemExists = struct{}{}

// New instantiates a new empty hset
func New() *Hset {
	return &Hset{items: make(map[interface{}]struct{})}
}

// Add adds the items (one or more) to the hset.
func (hset *Hset) Add(items ...interface{}) {
	hset.Lock()
	// defer hset.Unlock()
	for _, item := range items {
		hset.items[item] = itemExists
	}
	hset.Unlock()
}

// Remove removes the items (one or more) from the hset.
func (hset *Hset) Remove(items ...interface{}) {
	hset.Lock()
	// defer hset.Unlock()
	for _, item := range items {
		delete(hset.items, item)
	}
	hset.Unlock()
}

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

// Clear clears all values in the hset.
func (hset *Hset) Clear() {
	hset.Lock()
	// defer hset.Unlock()
	hset.items = make(map[interface{}]struct{})
	hset.Unlock()
}

// Exists returns a bool indicating if the given item exists in the set.
func (hset *Hset) Exists(item interface{}) bool {
	hset.RLock()
	_, ok := hset.items[item]
	hset.RUnlock()

	return ok
}

// Empty returns true if hset does not contain any elements.
func (hset *Hset) Empty() bool {
	// return hset.Size() == 0
	return hset.Len() == 0
}

// Size returns number of elements within the hset.
// func (hset *Hset) Size() int {
// 	return len(hset.items)
// }

// Len returns number of elements within the hset.
func (hset *Hset) Len() int {
	hset.RLock()
	size := len(hset.items)
	hset.RUnlock()
	return size
}

// Values returns all items in the hset.
// List()
func (hset *Hset) Values() []interface{} {
	hset.RLock()
	defer hset.RUnlock()

	// values := make([]interface{}, hset.Size())
	values := make([]interface{}, hset.Len())
	count := 0
	for item := range hset.items {
		values[count] = item
		count++
	}
	return values
}

// Same to determine whether the two hset type values are the same.
func (hset *Hset) Same(other Set) bool {
	hset.RLock()
	defer hset.RUnlock()

	if other == nil {
		return false
	}
	if hset.Len() != other.Len() {
		return false
	}

	for key := range hset.items {
		if !other.Contains(key) {
			return false
		}
	}
	return true
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
