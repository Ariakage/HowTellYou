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

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
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

	var db *sql.DB = openDataBase("root:qwerty123456@tcp(192.168.21.131:3306)/hty")

	var exec_res sql.Result
	//Create User Table
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_user ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT NOT NULL, `name` VARCHAR(16) NOT NULL, `email` VARCHAR(50) NOT NULL, `pwd` VARCHAR(20) NOT NULL, `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP );")
	fmt.Println(exec_res.RowsAffected())
	// Create Friend Table (https://blog.csdn.net/wo541075754/article/details/82733278)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_friend ( `user_id` INT NOT NULL, `friend_id` INT NOT NULL, `user_group` VARCHAR ( 10 ) NOT NULL, `friend_group` VARCHAR ( 10 ) NOT NULL );")
	fmt.Println(exec_res.RowsAffected())
	// Create Group Table (https://blog.csdn.net/php_xml/article/details/108690219)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_group ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT DEFAULT '', `name` VARCHAR ( 16 ) NOT NULL, `owner_id` INT NOT NULL, `admins` LONGTEXT NOT NULL DEFAULT '', `members` LONGTEXT NOT NULL, `type` INT NOT NULL, `remark` VARCHAR ( 200 ) NOT NULL DEFAULT '', `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL );")
	fmt.Println(exec_res.RowsAffected())
	// Create Message Table (https://blog.csdn.net/qq_42249896/article/details/104033697)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_message ( `id` INT PRIMARY KEY AUTO_INCREMENT, `send_user_id` INT NOT NULL, `receive_user_id` INT NOT NULL, `content` TEXT NOT NULL, `send_time` DATETIME NOT NULL );")
	fmt.Println(exec_res.RowsAffected())
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
