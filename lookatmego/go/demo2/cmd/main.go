package main

import (
	"fmt"
	"os"
	"plugin"

	"github.com/tigh-latte/lookatmego"
)

var input = []byte("name: john")

func main() {
	plug, err := plugin.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	symb, err := plug.Lookup("New")
	if err != nil {
		panic(err)
	}

	fn, ok := symb.(func() lookatmego.Plugin)
	if !ok {
		panic("yikes")
	}

	fmt.Println(fn().Render(input))
}
