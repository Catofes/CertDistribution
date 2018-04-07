package src

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

type Config struct {
	ListenOn  string
	StorePath string
}

// Load function load Config from file.
func (s *Config) Load(path string) *Config {
	s.ListenOn = "127.0.0.1:8080"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("[Fatal] Can not read Config file: ", err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		log.Fatal("[Fatal] Can not parse Config file. ", err)
	}
	return s
}
