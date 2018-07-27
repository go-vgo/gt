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
	"testing"

	"github.com/vcaesar/tt"
)

func TestSetAdd(t *testing.T) {
	set := New()
	set.Add()
	set.Add(1)
	set.Add(2)
	set.Add(2, 3)
	set.Add()

	actualValue := set.Empty()
	tt.False(t, actualValue)

	actualVal := set.Len()
	tt.Equal(t, 3, actualVal)
}

func TestSetContains(t *testing.T) {
	set := New()
	set.Add(3, 1, 2)
	set.Add(2, 3)
	set.Add()
	actualValue := set.Contains()
	tt.True(t, actualValue)

	actualValue = set.Contains(1)
	tt.True(t, actualValue)

	actualValue = set.Contains(1, 2, 3)
	tt.True(t, actualValue)

	actualValue = set.Contains(1, 2, 3, 4)
	tt.False(t, actualValue)
}

func TestSetRemove(t *testing.T) {
	set := New()
	set.Add(3, 1, 2)
	set.Remove()
	actualValue := set.Len()
	tt.Equal(t, 3, actualValue)

	set.Remove(1)
	actualValue = set.Len()
	tt.Equal(t, 2, actualValue)

	set.Remove(3)
	set.Remove(3)
	set.Remove()
	set.Remove(2)

	actualValue = set.Len()
	tt.Equal(t, 0, actualValue)
}
