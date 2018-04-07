package src

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"crypto/x509"
	"io/ioutil"
)

type certHandler struct {
	Config
	data *storage
}

func (s *certHandler) certPut(ctx iris.Context) {
	id := uuid.NewV4().String()
	file, _, err := ctx.FormFile("cert")
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	_, err = x509.ParseCertificate(data)
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	s.data.Set(id, data)
	ctx.StatusCode(200)
	ctx.WriteString(id)
}

func (s *certHandler) certPost(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	_, err := s.data.Get(id)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	file, _, err := ctx.FormFile("cert")
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	_, err = x509.ParseCertificate(data)
	if err != nil {
		ctx.StatusCode(400)
		return
	}
	s.data.Set(id, data)
	ctx.StatusCode(200)
}

func (s *certHandler) certGet(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	cert, err := s.data.Get(id)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	ctx.JSON(cert)
}

func (s *certHandler) certGetRaw(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	cert, err := s.data.Get(id)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	ctx.Write(cert.Raw)
}

func (s *certHandler) bind(app *iris.Application) {
	app.Get("/{cert_id: string}", s.certGet)
	app.Get("/{cert_id: string}/raw", s.certGetRaw)
	app.Put("/", s.certPut)
	app.Post("/{cert_id: string}", s.certPost)
}
