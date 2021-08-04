package main

import (
	"flag"
	"log"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	g := newGenerator()

	var flags flag.FlagSet

	// flags.StringVar(&g.OptionPrefix, "option_prefix", "sniper", "")
	// flags.StringVar(&g.RootPackage, "root_package", "sniper", "")
	// flags.BoolVar(&g.ValidateEnable, "validate_enable", false, "")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(g.Generate)

}
