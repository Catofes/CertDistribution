package main

import (
	"fmt"
	"flag"
	"github.com/Catofes/CertDistribution/src"
)

var _version_ string

func main() {
	fmt.Printf("Background Web. Version %s.\n", _version_)
	configPath := flag.String("c", "config.json", "Config file path.")
	flag.Parse()
	c := src.Config{}
	c.Load(*configPath)
	src.Run(&c)
}
