package pwd

import (
	"testing"

	"github.com/vcaesar/tt"
)

var (
	testPwd = "adc123"
)

func TestPwGen(t *testing.T) {
	pw := Gen(testPwd)
	b := Check(pw, testPwd)

	tt.True(t, b)
}
