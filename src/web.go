package src

import (
	"github.com/kataras/iris"
	"sync"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var (
	app  *iris.Application
	once sync.Once
)

func getApp() *iris.Application {
	once.Do(func() {
		app = iris.New()
		app.Configure(iris.WithRemoteAddrHeader("X-Real-Ip"),
			iris.WithRemoteAddrHeader("X-Forwarded-For"))
		app.Use(recover.New())
		app.Use(logger.New())
	})
	return app
}

func bind(c *Config) {
	aStorage := (&storage{Config: *c}).Init()
	aCertHandler := certHandler{*c, aStorage}
	aCertHandler.bind(app)
}

func Run(c *Config) {
	getApp()
	bind(c)
	app.Run(iris.Addr(c.ListenOn), iris.WithoutVersionChecker)
}
