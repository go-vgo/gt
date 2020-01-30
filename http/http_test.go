package http

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt"
)

func TestApi(t *testing.T) {
	m := Map{}
	r, e := Api("https://github.com/vcaesar/tt", m)
	fmt.Println("get: ", string(r))

	tt.Nil(t, e)
	tt.NotNil(t, r)
}
