package main

import (
	"fmt"

	"github.com/go-vgo/gt/file"
)

func main() {
	sha, err := file.OFileSha("../file/file.go", "sha256")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sha)

	filesize, err := file.FileSize("../file/flie.go")
	fmt.Println(filesize, err)
}
