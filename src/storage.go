package src

import (
	"crypto/x509"
	"encoding/json"
	"log"
	"io/ioutil"
	"sync"
	"errors"
	"encoding/pem"
)

type Cert struct {
	Chain []x509.Certificate
	Id    string
	Data  string
}

type storage struct {
	Config
	data  map[string]Cert
	mutex sync.Mutex
}

func (s *storage) Init() *storage {
	s.data = make(map[string]Cert)
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
	err = ioutil.WriteFile(s.StorePath, data, 500)
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
	tmp, err := s.ParseCert(cert)
	if err != nil {
		return err
	}
	s.data[Id] = Cert{tmp, Id, string(cert)}
	return nil
}

func (s *storage) Get(Id string) (Cert, error) {
	cert, ok := s.data[Id]
	if ok {
		return cert, nil
	} else {
		return cert, errors.New("cert not found")
	}
}

func (s *storage) ParseCert(cert []byte) ([]x509.Certificate, error) {
	chain := make([]x509.Certificate, 0)
	restPEMBlock := cert
	var certDERBlock *pem.Block = nil
	for {
		certDERBlock, restPEMBlock = pem.Decode(restPEMBlock)
		if certDERBlock == nil {
			return chain, nil
		}
		cert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			log.Println(err)
			continue
		}
		chain = append(chain, *cert)
	}
	if len(chain) <= 0 {
		return chain, errors.New("empty chain")
	}
	return chain, nil
}
