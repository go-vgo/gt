// +build !sync

package hset

import (
	"sync"
)

// Hset holds elements in go's native map
type Hset struct {
	items map[interface{}]struct{}
	// items sync.Map[interface{}]struct{}
	sync.RWMutex
}

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
