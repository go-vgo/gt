package file

import (
	"testing"

	"github.com/vcaesar/tt"
)

var testFile = "../testdata/file_test.txt"

func TestAppendTo(t *testing.T) {
	for index := 0; index < 10; index++ {
		tt.Nil(t, AppendTo(testFile, "test"))
	}

	// os.Truncate(testFile, 0)
	err := Empty(testFile)
	tt.Nil(t, err)

	r, err := Read(testFile)
	tt.Equal(t, "", r)
	tt.Nil(t, err)

	err = Write(testFile, "test")
	tt.Nil(t, err)

	r, err = Read(testFile)
	tt.Equal(t, "test", r)
	tt.Nil(t, err)
}

func TestList(t *testing.T) {
	f, err := List("./", ".go")
	tt.Nil(t, err)
	tt.Equal(t, 3, len(f))

	f, err = ListDir("./", ".go")
	tt.Nil(t, err)
	tt.Equal(t, 0, len(f))

	f, err = Walk("./", ".go")
	tt.Nil(t, err)
	tt.Equal(t, 3, len(f))

	f, err = WalkDir("./", ".go")
	tt.Nil(t, err)
	tt.Equal(t, 0, len(f))
}

func TestSys(t *testing.T) {
	h, e := Sha(testFile)
	tt.Equal(t, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", h)
	tt.Nil(t, e)

	s, e := Size(testFile)
	tt.Nil(t, e)
	tt.Equal(t, 4, s)
}
