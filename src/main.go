package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Hello, World!")
	})

	api := app.Party("/api")

	mvc.Configure(api.Party("/user"), userMvc)

	app.Run(iris.Addr(":8080"))
}
