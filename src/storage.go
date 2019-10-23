package src

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

type cert struct {
	Chain []x509.Certificate
	cert  *x509.Certificate
	ID    string
	Data  string
	mutex sync.Mutex
	c     context.Context
	cc    context.CancelFunc
}

func (s *cert) init() *cert {
	s.c, s.cc = context.WithCancel(context.Background())
	s.Chain = make([]x509.Certificate, 0)
	return s
}

func (s *cert) update(cert []byte) error {
	err := s.parseCert(cert)
	if err != nil {
		return err
	}
	s.cc()
	s.c, s.cc = context.WithCancel(context.Background())
	return nil
}

func (s *cert) parseCert(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	chain := make([]x509.Certificate, 0)
	restPEMBlock := data
	var certDERBlock *pem.Block
	for {
		certDERBlock, restPEMBlock = pem.Decode(restPEMBlock)
		if certDERBlock == nil {
			break
		}
		cert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			log.Println(err)
			continue
		}
		chain = append(chain, *cert)
	}
	if len(chain) <= 0 {
		return errors.New("empty chain")
	}
	var cert *x509.Certificate
	for k, v := range chain {
		if !v.IsCA {
			cert = &chain[k]
		}
	}
	if cert == nil {
		return errors.New("no terminate cert")
	}
	s.Chain = chain
	s.Data = string(data)
	s.cert = cert
	return nil
}

type storage struct {
	Config
	data  map[string]*cert
	mutex sync.Mutex
}

func (s *storage) init() *storage {
	s.data = make(map[string]*cert)
	return s
}

func (s *storage) save() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.data {
		v.mutex.Lock()
		defer v.mutex.Unlock()
	}
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

func (s *storage) load() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := ioutil.ReadFile(s.StorePath)
	if err != nil {
		log.Println("Load certificates failed.")
		return
	}
	json.Unmarshal(data, &(s.data))
	for _, v := range s.data {
		v.init()
		v.parseCert([]byte(v.Data))
	}
}

func (s *storage) set(id string, c []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_c := cert{ID: id, Data: string(c)}
	_c.init()
	err := _c.parseCert(c)
	if err != nil {
		return err
	}
	s.data[id] = &_c
	return nil
}

func (s *storage) get(id string) (*cert, error) {
	cert, ok := s.data[id]
	if ok {
		return cert, nil
	}
	return nil, errors.New("cert not found")
}
