package src

import (
	"crypto/x509"
	"encoding/json"
	"log"
	"io/ioutil"
	"sync"
	"errors"
)

type cert struct {
	x509.Certificate
	Id string
}

type storage struct {
	Config
	data  map[string]cert
	mutex sync.Mutex
}

func (s *storage) Init() *storage {
	s.data = make(map[string]cert)
	return s
}

func (s *storage) Save() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := json.Marshal(s.data)
	if err != nil {
		log.Println("Save certificates failed.")
		return
	}
	err = ioutil.WriteFile(s.StorePath, data, 700)
	if err != nil {
		log.Print("Save certificates failed.")
	}
}

func (s *storage) Load() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := ioutil.ReadFile(s.StorePath)
	if err != nil {
		log.Println("Load certificates failed.")
		return
	}
	json.Unmarshal(data, &(s.data))
}

func (s *storage) Set(Id string, cert []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tmp, err := x509.ParseCertificate(cert)
	if err != nil {
		return err
	}
	aCert := cert{*tmp, Id}
	s.data[Id] = aCert
	return nil
}

func (s *storage) Get(Id string) (cert, error) {
	cert, ok := s.data[Id]
	if ok {
		return cert, nil
	} else {
		return cert, errors.New("cert not found")
	}
}
