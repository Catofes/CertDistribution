package src

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
)

type certHandler struct {
	Config
	data *storage
}

func (s *certHandler) certPut(ctx iris.Context) {
	id := uuid.NewV4().String()
	file, _, err := ctx.FormFile("Cert")
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	_, err = s.data.ParseCert(fileData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	s.data.Set(id, fileData)
	s.data.Save()
	ctx.StatusCode(200)
	ctx.WriteString(id)
}

func (s *certHandler) certPost(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	_, err := s.data.Get(id)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(404)
		return
	}
	file, _, err := ctx.FormFile("Cert")
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	_, err = s.data.ParseCert(fileData)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(400)
		return
	}
	s.data.Set(id, fileData)
	s.data.Save()
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
	ctx.WriteString(cert.Data)
}

func (s *certHandler) bind(app *iris.Application) {
	app.Get("/{cert_id: string}", s.certGet)
	app.Get("/{cert_id: string}/raw", s.certGetRaw)
	app.Put("/", s.certPut)
	app.Post("/{cert_id: string}", s.certPost)
}
