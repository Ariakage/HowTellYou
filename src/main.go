/*
CREATED BY DR.ALANORAGE on 2023.07.10
*/
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Hello, World!")
	})

	api_prt := app.Party("/api")

	/**/
	user_prt := api_prt.Party("/user")

	user_loginAPI := user_prt.Party("/login")
	{

	}
	/**/

	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./config_iris.toml")))
}
