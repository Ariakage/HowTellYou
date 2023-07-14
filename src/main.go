/*
CREATED BY DR.ALANORAGE on 2023.07.10
*/
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
)

var port = flag.Int("p", 8080, "Server Listen Port")
var configFilePath = flag.String("cfgp", "./config_iris.toml", "Config File Path")

func main() {

	flag.Parse()
	if *port > 65536 || *port < 0 {
		fmt.Println("HTY Startup error - port parameter >65536 or <0")
		os.Exit(0)
	}

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	var db *sql.DB = openDataBase("./db.sqlite")
	//Create User Table
	execSQL(db, `
CREATE TABLE IF NOT EXISTS hty_user
(
id int primary key identity,
favimg text default '',
name varchar(16) not null,
email varchar(50) not null,
pwd varchar(20) not null,
create_time datetime DEFAULT CURRENT_TIMESTAMP not null
)
	`)
	// Create Friend Table (https://blog.csdn.net/wo541075754/article/details/82733278)
	execSQL(db, `
CREATE TABLE IF NOT EXISTS hty_friend
(
user_id int not null,
friend_id int not null,
user_group varchar(10) not null,
friend_group varchar(10) not null
)
		`)
	// Create Group Table (https://blog.csdn.net/php_xml/article/details/108690219)
	execSQL(db, `
CREATE TABLE IF NOT EXISTS hty_group
(
id int primary key identity,
favimg text default '',
name varchar(16) not null,
owner_id int not null,
admins longtext not null default '',
members longtext not null,
type int not null,
remark varchar(200) not null default '',
create_time datetime DEFAULT CURRENT_TIMESTAMP not null
)
	`)
	// Create Message Table (https://blog.csdn.net/qq_42249896/article/details/104033697)
	execSQL(db, `
CREATE TABLE IF NOT EXISTS hty_message
(
id int primary key identity,
send_user_id int not null,
receive_user_id int not null,
content text not null,
send_time datetime not null
)
	`)
	closeDB(db)

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

	app.Run(iris.Addr(":"+strconv.Itoa(*port)), iris.WithConfiguration(iris.TOML(*configFilePath)))
}
