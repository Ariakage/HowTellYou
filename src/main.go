/*
CREATED BY DR.ALANORAGE on 2023.07.10
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

var port = flag.Int("p", 8080, "Server Listen Port")

func main() {
	flag.Parse()
	if *port > 65536 || *port < 0 {
		fmt.Println("HTY Startup error - port parameter >65536 or <0")
		os.Exit(0)
	}

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Hello, World!")
	})

	api_prt := app.Party("/api")

	/* - User part - */
	user_prtAPI := api_prt.Party("/user")
	{
		user_prtAPI.Use(iris.Compression)

		//Login
		{
			user_prtAPI.Get("/login", func(ctx iris.Context) {
				ctx.Text("login - get")
			})
			user_prtAPI.Post("/login", func(ctx iris.Context) {
				ctx.Text("login - post")
			})
		}
		//Register
		{
			user_prtAPI.Get("/register", func(ctx iris.Context) {
				ctx.Text("register - get")
			})
			user_prtAPI.Post("/register", func(ctx iris.Context) {
				ctx.Text("register - post")
			})
		}
		//Logout
		{
			user_prtAPI.Get("/logout", func(ctx iris.Context) {
				ctx.Text("logout - get")
			})
			user_prtAPI.Post("/logout", func(ctx iris.Context) {
				ctx.Text("logout - post")
			})
		}
		//Find Account
		{
			user_prtAPI.Get("/find_account", func(ctx iris.Context) {
				ctx.Text("find-account - get")
			})
			user_prtAPI.Post("/find_account", func(ctx iris.Context) {
				ctx.Text("find-account - post")
			})
		}
	}
	/* --- */

	app.Run(iris.Addr(":"+strconv.Itoa(*port)), iris.WithConfiguration(iris.TOML("./config_iris.toml")))
}
