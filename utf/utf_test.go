package utf

import (
	"testing"

	"github.com/vcaesar/tt"
)

var (
	ucode = `\u4f60\u597d\"\\u4f60\\u597d\"`
	text  = `你好\"\你\好\"`
)

func TestUnicodeToUTF8(t *testing.T) {
	u := UnicodeToUTF8(ucode)
	tt.Equal(t, text, u)

	r := CharCodeAt(text, 1)
	tt.Equal(t, "22909", r)

	tt.Equal(t, `[U4f60 U597d \\ \" \\ U4f60 \\ U597d \\ \"]`, ToUC(text))
}
