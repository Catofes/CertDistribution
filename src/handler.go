package src

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

type certHandler struct {
	Config
	data *storage
}

func (s *certHandler) certPut(ctx iris.Context) {
	id := uuid.NewV4().String()
	file, _, err := ctx.FormFile("Cert")
	if err != nil {
		log.Printf("Get post cert failed: %s.", err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file failed: %s.", err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	err = s.data.set(id, fileData)
	if err != nil {
		log.Printf("Save cert failed: %s.", err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	s.data.save()
	ctx.StatusCode(200)
	ctx.WriteString(id)
}

func (s *certHandler) certWait(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	c, err := s.data.Get(id)
	if err != nil {
		log.Println(err)
		ctx.Text(err.Error())
		ctx.StatusCode(404)
		return
	}
	serialID := ctx.Params().Get("sn")
	if serialID != c.cert.SerialNumber.String() {
		ctx.Text(c.Data)
	} else {
		select {
		case <-c.c.Done():
			ctx.Text(c.Data)
		case <-time.After(1 * time.Hour):
			ctx.StatusCode(204)
		}
	}
}

func (s *certHandler) certPost(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	c, err := s.data.Get(id)
	if err != nil {
		log.Println(err)
		ctx.Text(err.Error())
		ctx.StatusCode(404)
		return
	}
	file, _, err := ctx.FormFile("Cert")
	if err != nil {
		log.Println(err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	defer file.Close()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	err = c.update(fileData)
	if err != nil {
		log.Println(err)
		ctx.Text(err.Error())
		ctx.StatusCode(400)
		return
	}
	ctx.StatusCode(200)
}

func (s *certHandler) certGet(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	cert, err := s.data.Get(id)
	if err != nil {
		ctx.Text(err.Error())
		ctx.StatusCode(404)
		return
	}
	ctx.JSON(cert)
}

func (s *certHandler) certGetRaw(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	cert, err := s.data.Get(id)
	if err != nil {
		ctx.Text(err.Error())
		ctx.StatusCode(404)
		return
	}
	ctx.WriteString(cert.Data)
}

func (s *certHandler) certSN(ctx iris.Context) {
	id := ctx.Params().Get("cert_id")
	cert, err := s.data.Get(id)
	if err != nil {
		ctx.Text(err.Error())
		ctx.StatusCode(404)
		return
	}
	ctx.Text(cert.cert.SerialNumber.String())
}

func (s *certHandler) bind(app *iris.Application) {
	app.Get("/{cert_id: string}", s.certGet)
	app.Get("/{cert_id: string}/raw", s.certGetRaw)
	app.Get("/{cert_id: string}/wait/{sn: string}", s.certWait)
	app.Get("/{cert_id: string}/sn", s.certSN)
	app.Put("/", s.certPut)
	app.Post("/{cert_id: string}", s.certPost)
}
