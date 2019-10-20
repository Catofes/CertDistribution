package main

import (
	"flag"
	"fmt"

	"github.com/Catofes/CertDistribution/src"
)

var _version string

func main() {
	fmt.Printf("Background Web. Version %s.\n", _version)
	configPath := flag.String("c", "config.json", "Config file path.")
	flag.Parse()
	c := src.Config{}
	c.Load(*configPath)
	src.Run(&c)
}
