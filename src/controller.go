package main

import (
	//"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

/**/
func userMvc(app *mvc.Application) {
	app.Handle(new(UserController))
}

type UserController struct{}

/**/
