package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	g := markdown{}

	var flags flag.FlagSet

	flags.StringVar(&g.Prefix, "prefix", "/", "API path prefix")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(g.Generate)
}
