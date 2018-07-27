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
	Empty(testFile)

	r, err := Read(testFile)
	tt.Equal(t, "", r)
	tt.Nil(t, err)

	Write(testFile, "test")
	r, err = Read(testFile)
	tt.Equal(t, "test", r)
	tt.Nil(t, err)
}
