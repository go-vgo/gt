package hset

import (
	"testing"

	"github.com/vcaesar/tt"
)

var set = New()

func BenchmarkAdd(b *testing.B) {
	fn := func() {
		set.Add()
		set.Add(1)
		set.Add(2)
		set.Add(2, 3)
		set.Add()
	}

	tt.BM(b, fn)
}

func BenchmarkRemove(b *testing.B) {
	fn := func() {
		set.Remove(3)
		set.Remove(3)
		set.Remove()
		set.Remove(2)
	}

	tt.BM(b, fn)
}
