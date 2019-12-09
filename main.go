package main

import (
	"fmt"

	"github.com/williamlsh/goexamples/build"
)

var Version = "development"

func main() {
	fmt.Println("Version:\t", Version)
	fmt.Println("build.Time:\t", build.Time)
	fmt.Println("build.User:\t", build.User)
}
