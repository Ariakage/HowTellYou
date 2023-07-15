/*
CREATED BY DR.ALANORAGE on 2023.07.10
*/
package main

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/pelletier/go-toml"
)

var configFilePath = flag.String("cfgp", "./config.toml", "Config File Path")
var globalTokenMap map[int]*list.List = make(map[int]*list.List)

func main() {
	flag.Parse()

	config, _ := toml.LoadFile(*configFilePath)

	sv_cfgp := config.Get("server.config_path").(string)
	sv_port := int(config.Get("server.port").(int64))

	db_ip := config.Get("database.ip").(string)
	db_port := int(config.Get("database.port").(int64))
	db_user := config.Get("database.user").(string)
	db_pwd := config.Get("database.password").(string)
	db_dbname := config.Get("database.database_name").(string)
	db_cfg := mysql.Config{
		User:                 db_user,
		Passwd:               db_pwd,
		Net:                  "tcp",
		Addr:                 db_ip + ":" + strconv.Itoa(db_port),
		DBName:               db_dbname,
		AllowNativePasswords: true,
	}

	se_chars := config.Get("secure.allowed_chars").(string)
	se_saltlength := int(config.Get("secure.salt_length").(int64))
	se_key := config.Get("secure.key").(string)

	if sv_port > 65536 || sv_port < 0 {
		fmt.Println("HTY Startup error - port parameter >65536 or <0")
		os.Exit(0)
	}

	/*Debug Code There*/
	//fmt.Println(generateRandomString(10, config.Get("secure.allowed_chars").(string)))
	//fmt.Println(generateLoginToken(config.Get("secure.allowed_chars").(string), int(config.Get("secure.salt_length").(int64)), config.Get("secure.key").(string), 114514, time.Now()))
	/* --- */

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	//var db *sql.DB = openDataBase(db_link) //"root:qwerty123456@tcp(192.168.21.131:3306)/hty"
	var db *sql.DB = openDataBase(db_cfg.FormatDSN())
	//var db *sql.DB = openDataBase("root:qwerty123456@tcp(192.168.21.131:3306)/hty")

	/*var exec_res sql.Result
	// Create User Table
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_user ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT NOT NULL, `name` VARCHAR(16) NOT NULL, `email` VARCHAR(50) NOT NULL, `pwd` VARCHAR(512) NOT NULL, `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP );")
	fmt.Println(exec_res.RowsAffected())
	// Create Friend Table (https://blog.csdn.net/wo541075754/article/details/82733278)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_friend ( `user_id` INT NOT NULL, `friend_id` INT NOT NULL, `user_group` VARCHAR ( 10 ) NOT NULL, `friend_group` VARCHAR ( 10 ) NOT NULL );")
	fmt.Println(exec_res.RowsAffected())
	// Create Group Table (https://blog.csdn.net/php_xml/article/details/108690219)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_group ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT DEFAULT '', `name` VARCHAR ( 16 ) NOT NULL, `owner_id` INT NOT NULL, `admins` LONGTEXT NOT NULL DEFAULT '', `members` LONGTEXT NOT NULL, `type` INT NOT NULL, `remark` VARCHAR ( 200 ) NOT NULL DEFAULT '', `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL );")
	fmt.Println(exec_res.RowsAffected())
	// Create Message Table (https://blog.csdn.net/qq_42249896/article/details/104033697)
	exec_res = execSQL(db, "CREATE TABLE IF NOT EXISTS hty_message ( `id` INT PRIMARY KEY AUTO_INCREMENT, `send_user_id` INT NOT NULL, `receive_user_id` INT NOT NULL, `content` TEXT NOT NULL, `send_time` DATETIME NOT NULL );")
	fmt.Println(exec_res.RowsAffected())*/

	// Create Tables to mysql db
	createTables(db)

	//closeDB(db)

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>How Tell You Back</title>
		<style>
			*{
				text-align: center;
			}
		</style>
	</head>
	<body>
		<h1>This is How Tell You Project Back Side</h1>
		<h1>Please do not directly access this site</h1>
		<h1>Thanks a lot!</h1>
	</body>
</html>
`)
	})

	api_prt := app.Party("/api")

	/* - User part - */
	user_prtAPI := api_prt.Party("/user")
	{
		user_prtAPI.Use(iris.Compression)

		//Login
		{
			user_prtAPI.Get("/login", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/login", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				name := respone_content["name"]
				pwd := respone_content["pwd"]
				if name == nil || pwd == nil {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "ID or PWD is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				row := db.QueryRow("SELECT `pwd` FROM hty_user WHERE `id` = ? OR `name` = ?", name.(string), name.(string))
				var pwd_dbd string
				err := row.Scan(&pwd_dbd)
				if err != nil {
					fmt.Println(err)
				}
				if pwd_dbd == GetSHA256HashCode([]byte(pwd.(string))) {
					row = db.QueryRow("SELECT `id` FROM hty_user WHERE `id` = ? OR `name` = ?", name.(string), name.(string))
					var id_dbd int
					err := row.Scan(&id_dbd)
					if err != nil {
						fmt.Println(err)
					}
					tk, slt := GenerateLoginToken(se_chars, se_saltlength, se_key, id_dbd, time.Now())
					_ = slt
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "success"
					d["token"] = tk
					var m *map[string]interface{} = makeResponse(200, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id_dbd]; ok {
						//globalTokenMap[int(id.(float64))][len(globalTokenMap[int(id.(float64))])] = tk
						var mp map[string]string = make(map[string]string)
						mp["token"] = tk
						mp["salt"] = slt
						globalTokenMap[id_dbd].PushBack(mp)

					} else {
						//var sarr []string = make([]string, 1)
						//sarr[0] = tk
						//globalTokenMap[int(id.(float64))] = sarr
						globalTokenMap[id_dbd] = list.New()
						var mp map[string]string = make(map[string]string)
						mp["token"] = tk
						mp["salt"] = slt
						globalTokenMap[id_dbd].PushBack(mp)
					}
					//fmt.Println(globalTokenMap[int(id.(float64))].Back().Value.(map[string]string))
					ctx.Text(string(m_b))
				} else {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "ID or PWD err"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
			})
		}
		//Register
		{
			user_prtAPI.Get("/register", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/register", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				name := respone_content["name"]
				email := respone_content["email"]
				pwd := respone_content["pwd"]
				if name == nil || email == nil || pwd == nil {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "name, email or pwd is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				res, err := db.Exec("INSERT INTO hty_user(`favimg` ,`name`, `nickname`, `email`, `pwd`) VALUES (?, ?, ?, ?, ?)", "", name.(string), name.(string), email.(string), GetSHA256HashCode([]byte(pwd.(string))))
				if err != nil {
					fmt.Println(err)
				}
				_ = res
				var d map[string]interface{} = make(map[string]interface{})
				d["status"] = "success"
				var m *map[string]interface{} = makeResponse(0, d)
				m_b, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err)
				}
				ctx.Text(string(m_b))
			})
		}
		//Logout
		{
			user_prtAPI.Get("/logout", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/logout", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				name := respone_content["name"]
				tk := respone_content["token"]
				if name == nil || tk == nil {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "Name, Id or Token is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				row := db.QueryRow("SELECT `id` FROM hty_user WHERE `id` = ? OR `name` = ?", name.(string), name.(string))
				var id_dbd int
				err := row.Scan(&id_dbd)
				if err != nil {
					fmt.Println(err)
				}
				if _, ok := globalTokenMap[id_dbd]; ok {
					for e := globalTokenMap[id_dbd].Front(); e != nil; e = e.Next() {
						if e.Value.(map[string]string)["token"] == tk {
							globalTokenMap[id_dbd].Remove(e)
							var d map[string]interface{} = make(map[string]interface{})
							d["status"] = "success"
							var m *map[string]interface{} = makeResponse(200, d)
							m_b, err := json.Marshal(m)
							if err != nil {
								fmt.Println(err)
							}
							ctx.Text(string(m_b))
							return
						}
					}
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "token is invalid"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))

				} else {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "id or name is invalid"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
			})
		}
		//Find Account
		{
			user_prtAPI.Get("/find_account", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/find_account", func(ctx iris.Context) {
				ctx.Text("find-account - post")
			})
		}
	}
	/* --- */

	app.Run(iris.Addr(":"+strconv.Itoa(sv_port)), iris.WithConfiguration(iris.TOML(sv_cfgp)))
}

/* Private */
func createTables(db *sql.DB) {
	//Create User Table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS hty_user ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT NOT NULL, `name` VARCHAR(16) NOT NULL, `nickname` VARCHAR (20) NOT NULL, `email` VARCHAR(50) NOT NULL, `pwd` VARCHAR(512) NOT NULL, `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP );")
	if err != nil {
		fmt.Println(err)
	}

	// Create Friend Table (https://blog.csdn.net/wo541075754/article/details/82733278)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS hty_friend ( `user_id` INT NOT NULL, `friend_id` INT NOT NULL, `user_group` VARCHAR ( 10 ) NOT NULL, `friend_group` VARCHAR ( 10 ) NOT NULL );")
	if err != nil {
		fmt.Println(err)
	}
	// Create Group Table (https://blog.csdn.net/php_xml/article/details/108690219)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS hty_group ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT DEFAULT '', `name` VARCHAR ( 16 ) NOT NULL, `owner_id` INT NOT NULL, `admins` LONGTEXT NOT NULL DEFAULT '', `members` LONGTEXT NOT NULL, `type` INT NOT NULL, `remark` VARCHAR ( 200 ) NOT NULL DEFAULT '', `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL );")
	if err != nil {
		fmt.Println(err)
	}
	// Create Message Table (https://blog.csdn.net/qq_42249896/article/details/104033697)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS hty_message ( `id` INT PRIMARY KEY AUTO_INCREMENT, `send_user_id` INT NOT NULL, `receive_user_id` INT NOT NULL, `content` TEXT NOT NULL, `send_time` DATETIME NOT NULL );")
	if err != nil {
		fmt.Println(err)
	}
}

/* --- */

func GetRequestParams(ctx iris.Context) interface{} {
	var params map[string]interface{}
	_ = ctx.ReadJSON(&params)
	return &params
}

func makeResponse(status_code int, data interface{}) *map[string]interface{} {
	var m map[string]interface{} = make(map[string]interface{})
	m["code"] = status_code
	m["result"] = data
	return &m
}
