package conf

import (
	_ "embed"
	"log"
	"testing"

	"github.com/vcaesar/tt"
)

var (
	//go:embed conf.toml
	conf1 string
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

func TestEmbed(t *testing.T) {
	toml := Toml{}
	err := Init(conf1, &toml, true)

	log.Println("toml: ", toml)
	tt.Nil(t, err)
	tt.Equal(t, "conf", toml.Test)
}
