package conf

import (
	"log"
	"testing"

	"github.com/vcaesar/tt"
)

type Toml struct {
	Test string `toml:"test"`
}

func TestToml(t *testing.T) {
	toml := Toml{}
	err := Init("../testdata/conf.toml", &toml)

	log.Println("toml: ", toml)
	tt.Nil(t, err)
	tt.Equal(t, "conf", toml.Test)
}
