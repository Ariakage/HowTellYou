/*
CREATED BY DR.ALANORAGE on 2023.07.10
*/
package main

import (
	"container/list"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/pelletier/go-toml"
)

var configFilePath = flag.String("cfgp", "./config.toml", "Config File Path")
var globalTokenMap map[int]*list.List = make(map[int]*list.List)
var globalFindAccountCodeMap map[string]*list.List = make(map[string]*list.List)

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
		ParseTime:            true,
	}

	se_chars := config.Get("secure.allowed_chars").(string)
	se_saltlength := int(config.Get("secure.salt_length").(int64))
	se_key := config.Get("secure.key").(string)
	se_vc_tom := int(config.Get("secure.verifcode_timeout").(int64))

	ml_saddr := config.Get("mail.send_addr").(string)
	ml_sv := config.Get("mail.server").(string)
	ml_port := int(config.Get("mail.port").(int64))
	ml_ac := config.Get("mail.authcode").(string)
	ml_findaccount_html := config.Get("mail.find_account_html").(string)

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
				row := db.QueryRow("SELECT COUNT(*) FROM hty_user WHERE `name` = ?", name.(string))
				var ct_dbd int
				err := row.Scan(&ct_dbd)
				if err != nil {
					fmt.Println(err)
				}
				if ct_dbd == 0 {
					_, err := db.Exec("INSERT INTO hty_user(`favimg` ,`name`, `nickname`, `gender`, `description`, `email`, `pwd`) VALUES (?, ?, ?, ?, ?, ?, ?)", "", name.(string), name.(string), -1, "", email.(string), GetSHA256HashCode([]byte(pwd.(string))))
					if err != nil {
						fmt.Println(err)
					}
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "success"
					var m *map[string]interface{} = makeResponse(200, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				} else {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "name already exists"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
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
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				email := respone_content["email"]
				name := respone_content["name"]
				if email == nil || name == nil {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "Name, or Email is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				row := db.QueryRow("SELECT COUNT(*) FROM hty_user WHERE `name` = ? AND `email` = ?", name.(string), email.(string))
				var ct_dbd int
				err := row.Scan(&ct_dbd)
				if err != nil {
					fmt.Println(err)
				}
				if ct_dbd == 0 {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "user not found"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				} else {
					code := strconv.Itoa(GenerateRealRand(9)) + strconv.Itoa(GenerateRealRand(9)) + strconv.Itoa(GenerateRealRand(9)) +
						strconv.Itoa(GenerateRealRand(9)) + strconv.Itoa(GenerateRealRand(9)) + strconv.Itoa(GenerateRealRand(9))
					hl := strings.Replace(strings.Replace(ml_findaccount_html, "$CODE", code, -1), "$TIME", time.Now().String(), -1)
					SendEmail(ml_saddr, email.(string), ml_sv, ml_port, ml_ac, hl)
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "success"
					row = db.QueryRow("SELECT `id` FROM hty_user WHERE `name` = ? AND `email` = ?", name.(string), email.(string))
					var id_dbd int
					err := row.Scan(&id_dbd)
					if err != nil {
						fmt.Println(err)
					}
					tk, slt := GenerateFindAccountToken(se_chars, se_saltlength, se_key, id_dbd, code, time.Now())
					d["token"] = tk
					var m *map[string]interface{} = makeResponse(200, d)
					m_b, err := json.Marshal(m)
					if _, ok := globalFindAccountCodeMap[code]; ok {
						var mp map[string]string = make(map[string]string)
						mp["token"] = tk
						mp["salt"] = slt
						globalFindAccountCodeMap[code].PushBack(mp)
					} else {
						globalFindAccountCodeMap[code] = list.New()
						var mp map[string]string = make(map[string]string)
						mp["token"] = tk
						mp["salt"] = slt
						globalFindAccountCodeMap[code].PushBack(mp)
					}
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
				//SendEmail(ml_saddr, "2026159790@qq.com", ml_sv, ml_port, ml_ac, ml_findaccount_html)
			})

			user_prtAPI.Get("/find_account_code", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/find_account_code", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				pwd := respone_content["pwd"]
				code := respone_content["code"]
				tk := respone_content["token"]
				if code == nil || tk == nil || pwd == nil {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "code, pwd or token is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				if _, ok := globalFindAccountCodeMap[code.(string)]; ok {
					for e := globalFindAccountCodeMap[code.(string)].Front(); e != nil; e = e.Next() {
						if e.Value.(map[string]string)["token"] == tk {
							globalFindAccountCodeMap[code.(string)].Remove(e)
							t, err := SCDecryptString(e.Value.(map[string]string)["token"], se_key, "aes")
							if err != nil {
								fmt.Println(err)
							}
							r, _ := base64.StdEncoding.DecodeString(t)
							t1, _ := time.Parse(strings.Split(strings.Split(string(r), "@")[1], "$")[0], strings.Split(strings.Split(string(r), "@")[1], "$")[0])
							t2 := time.Now()
							rt := t1.Sub(t2)
							if rt.Minutes() <= float64(se_vc_tom) {
								uid := strings.Split(string(r), "@")[0]
								pwds := GetSHA256HashCode([]byte(pwd.(string)))
								_, err = db.Exec("UPDATE hty_user SET `pwd`=? WHERE `id` = ?", pwds, uid)
								if err != nil {
									fmt.Println(err)
								}

								var d map[string]interface{} = make(map[string]interface{})
								d["status"] = "success"
								var m *map[string]interface{} = makeResponse(200, d)
								m_b, err := json.Marshal(m)
								if err != nil {
									fmt.Println(err)
								}
								ctx.Text(string(m_b))
							} else {
								var d map[string]interface{} = make(map[string]interface{})
								d["status"] = "failed"
								d["reason"] = "code timeout"
								var m *map[string]interface{} = makeResponse(0, d)
								m_b, err := json.Marshal(m)
								if err != nil {
									fmt.Println(err)
								}
								ctx.Text(string(m_b))
							}
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
					d["reason"] = "code is invalid"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
			})
		}
		//Modify User Info
		{
			user_prtAPI.Get("/modify", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/modify", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				switch ctx.URLParam("type") {
				case "favimg":
					tk := respone_content["token"]
					img := respone_content["img"]
					if tk == nil || img == nil {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token or img is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					t, err := SCDecryptString(tk.(string), se_key, "aes")
					if err != nil {
						fmt.Println(err)
					}
					r, _ := base64.StdEncoding.DecodeString(t)
					uid := strings.Split(string(r), "@")[0]
					id, err := strconv.Atoi(uid)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id]; ok {
						for e := globalTokenMap[id].Front(); e != nil; e = e.Next() {
							if e.Value.(map[string]string)["token"] == tk {
								_, err := db.Exec("UPDATE hty_user SET `favimg`=? WHERE `id` = ?", img.(string), id)
								if err != nil {
									fmt.Println(err)
								}
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
						return
					} else {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
				case "nickname":
					tk := respone_content["token"]
					nickname := respone_content["nickname"]
					if tk == nil || nickname == nil {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token or nickname is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					t, err := SCDecryptString(tk.(string), se_key, "aes")
					if err != nil {
						fmt.Println(err)
					}
					r, _ := base64.StdEncoding.DecodeString(t)
					uid := strings.Split(string(r), "@")[0]
					id, err := strconv.Atoi(uid)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id]; ok {
						for e := globalTokenMap[id].Front(); e != nil; e = e.Next() {
							if e.Value.(map[string]string)["token"] == tk {
								_, err := db.Exec("UPDATE hty_user SET `nickname`=? WHERE `id` = ?", nickname.(string), id)
								if err != nil {
									fmt.Println(err)
								}
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
						return
					} else {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
				case "gender":
					tk := respone_content["token"]
					gender := respone_content["gender"]
					if tk == nil || gender == nil {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token or gender is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					if gender != -1 && gender != 0 && gender != 1 && gender != 2 {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "gender is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					t, err := SCDecryptString(tk.(string), se_key, "aes")
					if err != nil {
						fmt.Println(err)
					}
					r, _ := base64.StdEncoding.DecodeString(t)
					uid := strings.Split(string(r), "@")[0]
					id, err := strconv.Atoi(uid)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id]; ok {
						for e := globalTokenMap[id].Front(); e != nil; e = e.Next() {
							if e.Value.(map[string]string)["token"] == tk {
								_, err := db.Exec("UPDATE hty_user SET `gender`=? WHERE `id` = ?", int(gender.(float64)), id)
								if err != nil {
									fmt.Println(err)
								}
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
						return
					} else {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
				case "description":
					tk := respone_content["token"]
					desc := respone_content["description"]
					if tk == nil || desc == nil {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token or gender is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					t, err := SCDecryptString(tk.(string), se_key, "aes")
					if err != nil {
						fmt.Println(err)
					}
					r, _ := base64.StdEncoding.DecodeString(t)
					uid := strings.Split(string(r), "@")[0]
					id, err := strconv.Atoi(uid)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id]; ok {
						for e := globalTokenMap[id].Front(); e != nil; e = e.Next() {
							if e.Value.(map[string]string)["token"] == tk {
								_, err := db.Exec("UPDATE hty_user SET `description`=? WHERE `id` = ?", desc.(string), id)
								if err != nil {
									fmt.Println(err)
								}
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
						return
					} else {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
				case "pwd":
					tk := respone_content["token"]
					pwd := respone_content["pwd"]
					if tk == nil || pwd == nil {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token or pwd is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					t, err := SCDecryptString(tk.(string), se_key, "aes")
					if err != nil {
						fmt.Println(err)
					}
					r, _ := base64.StdEncoding.DecodeString(t)
					uid := strings.Split(string(r), "@")[0]
					id, err := strconv.Atoi(uid)
					if err != nil {
						fmt.Println(err)
					}
					if _, ok := globalTokenMap[id]; ok {
						for e := globalTokenMap[id].Front(); e != nil; e = e.Next() {
							if e.Value.(map[string]string)["token"] == tk {
								delete(globalTokenMap, id)
								_, err := db.Exec("UPDATE hty_user SET `pwd`=? WHERE `id` = ?", GetSHA256HashCode([]byte(pwd.(string))), id)
								if err != nil {
									fmt.Println(err)
								}
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
						return
					} else {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "token is invalid"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
				default:
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "type is invalid"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
			})
		}
		//Query User Info
		{
			user_prtAPI.Get("/query", func(ctx iris.Context) {
				id := ctx.URLParam("id")
				name := ctx.URLParam("name")
				if id == "" && name == "" {
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "id and name is null"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
					return
				}
				var d map[string]interface{} = make(map[string]interface{})
				var row *sql.Row
				if id == "" && name != "" {
					row = db.QueryRow("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `name` = ?", name)
				} else if id != "" && name == "" {
					row = db.QueryRow("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `id` = ?", id)
				} else {
					row = db.QueryRow("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `name` = ? AND `id` = ?", name, id)
				}
				var (
					id_          int
					favimg_      string
					name_        string
					nickname_    string
					gender_      int
					description_ string
					email_       string
					create_time_ time.Time
				)
				err := row.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
				if err != nil {
					fmt.Println(err)
				}
				d["id"] = id_
				d["favimg"] = favimg_
				d["name"] = name_
				d["nickname"] = nickname_
				d["gender"] = gender_
				d["description"] = description_
				d["email"] = email_
				d["create_time"] = create_time_
				var dx map[string]interface{} = make(map[string]interface{})
				dx["status"] = "success"
				dx["data"] = d
				var m *map[string]interface{} = makeResponse(200, dx)
				m_b, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err)
				}
				ctx.Text(string(m_b))
			})
			user_prtAPI.Post("/query", func(ctx iris.Context) {
				ctx.Text("Post request is not supported")
			})
		}
		//Search User
		{
			user_prtAPI.Get("/search", func(ctx iris.Context) {
				ctx.Text("Get request is not supported")
			})
			user_prtAPI.Post("/search", func(ctx iris.Context) {
				respone_content := *(GetRequestParams(ctx).(*map[string]interface{}))
				switch ctx.URLParam("type") {
				case "id":
					id := respone_content["id"]
					if id == "" {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "id is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					var d map[string]interface{} = make(map[string]interface{})
					rows, err := db.Query("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `id` = ?", id)
					if err != nil {
						fmt.Println(err)
					}
					lst := list.New()
					for rows.Next() {
						var (
							id_          int
							favimg_      string
							name_        string
							nickname_    string
							gender_      int
							description_ string
							email_       string
							create_time_ time.Time
						)
						err := rows.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
						var mp map[string]interface{} = make(map[string]interface{})
						mp["id"] = id_
						mp["favimg"] = favimg_
						mp["name"] = name_
						mp["nickname"] = nickname_
						mp["gender"] = gender_
						mp["description"] = description_
						mp["email"] = email_
						mp["create_time"] = create_time_
						lst.PushBack(mp)
						fmt.Println(lst)
						if err != nil {
							fmt.Println(err)
						}
					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}
					rows.Close()
					arr := make([]interface{}, 0)
					for i := lst.Front(); i != nil; i = i.Next() {
						arr = append(arr, i.Value)
					}
					d["status"] = "success"
					d["data"] = arr
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				case "name":
					name := respone_content["name"]
					if name == "" {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "name is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					var d map[string]interface{} = make(map[string]interface{})
					rows, err := db.Query("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `name` = ?", name)
					if err != nil {
						fmt.Println(err)
					}
					lst := list.New()
					for rows.Next() {
						var (
							id_          int
							favimg_      string
							name_        string
							nickname_    string
							gender_      int
							description_ string
							email_       string
							create_time_ time.Time
						)
						err := rows.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
						var mp map[string]interface{} = make(map[string]interface{})
						mp["id"] = id_
						mp["favimg"] = favimg_
						mp["name"] = name_
						mp["nickname"] = nickname_
						mp["gender"] = gender_
						mp["description"] = description_
						mp["email"] = email_
						mp["create_time"] = create_time_
						lst.PushBack(mp)
						if err != nil {
							fmt.Println(err)
						}
					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}
					rows.Close()
					arr := make([]interface{}, 0)
					for i := lst.Front(); i != nil; i = i.Next() {
						arr = append(arr, i.Value)
					}
					d["status"] = "success"
					d["data"] = arr
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				case "nickname":
					nickname := respone_content["nickname"]
					if nickname == "" {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "nickname is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					var d map[string]interface{} = make(map[string]interface{})
					rows, err := db.Query("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `nickname` = ?", nickname)
					if err != nil {
						fmt.Println(err)
					}
					lst := list.New()
					for rows.Next() {
						var (
							id_          int
							favimg_      string
							name_        string
							nickname_    string
							gender_      int
							description_ string
							email_       string
							create_time_ time.Time
						)
						err := rows.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
						var mp map[string]interface{} = make(map[string]interface{})
						mp["id"] = id_
						mp["favimg"] = favimg_
						mp["name"] = name_
						mp["nickname"] = nickname_
						mp["gender"] = gender_
						mp["description"] = description_
						mp["email"] = email_
						mp["create_time"] = create_time_
						lst.PushBack(mp)
						if err != nil {
							fmt.Println(err)
						}
					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}
					rows.Close()
					arr := make([]interface{}, 0)
					for i := lst.Front(); i != nil; i = i.Next() {
						arr = append(arr, i.Value)
					}
					d["status"] = "success"
					d["data"] = arr
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				case "gender":
					gender := respone_content["gender"]
					if gender == "" {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "gender is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					var d map[string]interface{} = make(map[string]interface{})
					rows, err := db.Query("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `id` = ?", gender)
					if err != nil {
						fmt.Println(err)
					}
					lst := list.New()
					for rows.Next() {
						var (
							id_          int
							favimg_      string
							name_        string
							nickname_    string
							gender_      int
							description_ string
							email_       string
							create_time_ time.Time
						)
						err := rows.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
						var mp map[string]interface{} = make(map[string]interface{})
						mp["id"] = id_
						mp["favimg"] = favimg_
						mp["name"] = name_
						mp["nickname"] = nickname_
						mp["gender"] = gender_
						mp["description"] = description_
						mp["email"] = email_
						mp["create_time"] = create_time_
						lst.PushBack(mp)
						if err != nil {
							fmt.Println(err)
						}
					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}
					rows.Close()
					arr := make([]interface{}, 0)
					for i := lst.Front(); i != nil; i = i.Next() {
						arr = append(arr, i.Value)
					}
					d["status"] = "success"
					d["data"] = arr
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				case "email":
					email := respone_content["email"]
					if email == "" {
						var d map[string]interface{} = make(map[string]interface{})
						d["status"] = "failed"
						d["reason"] = "email is null"
						var m *map[string]interface{} = makeResponse(0, d)
						m_b, err := json.Marshal(m)
						if err != nil {
							fmt.Println(err)
						}
						ctx.Text(string(m_b))
						return
					}
					var d map[string]interface{} = make(map[string]interface{})
					rows, err := db.Query("SELECT `id`, `favimg`, `name`, `nickname`, `gender`, `description`, `email`, `create_time` FROM hty_user WHERE `email` = ?", email)
					if err != nil {
						fmt.Println(err)
					}
					lst := list.New()
					for rows.Next() {
						var (
							id_          int
							favimg_      string
							name_        string
							nickname_    string
							gender_      int
							description_ string
							email_       string
							create_time_ time.Time
						)
						err := rows.Scan(&id_, &favimg_, &name_, &nickname_, &gender_, &description_, &email_, &create_time_)
						var mp map[string]interface{} = make(map[string]interface{})
						mp["id"] = id_
						mp["favimg"] = favimg_
						mp["name"] = name_
						mp["nickname"] = nickname_
						mp["gender"] = gender_
						mp["description"] = description_
						mp["email"] = email_
						mp["create_time"] = create_time_
						lst.PushBack(mp)
						if err != nil {
							fmt.Println(err)
						}
					}
					err = rows.Err()
					if err != nil {
						fmt.Println(err)
					}
					rows.Close()
					arr := make([]interface{}, 0)
					for i := lst.Front(); i != nil; i = i.Next() {
						arr = append(arr, i.Value)
					}
					d["status"] = "success"
					d["data"] = arr
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				case "description":
					description := respone_content["description"]
					_ = description
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "this function can't work because it is didnt be done"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				default:
					var d map[string]interface{} = make(map[string]interface{})
					d["status"] = "failed"
					d["reason"] = "type is invalid"
					var m *map[string]interface{} = makeResponse(0, d)
					m_b, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
					}
					ctx.Text(string(m_b))
				}
			})
		}
	}
	/* --- */
	/* Group Part */
	/* --- */
	/* Message Part */
	/* --- */
	/* Other Part */
	/* --- */

	app.Run(iris.Addr(":"+strconv.Itoa(sv_port)), iris.WithConfiguration(iris.TOML(sv_cfgp)))
}

/* Private */
func createTables(db *sql.DB) {
	//Create User Table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS hty_user ( `id` INT PRIMARY KEY AUTO_INCREMENT, `favimg` TEXT NOT NULL, `name` VARCHAR(16) UNIQUE NOT NULL, `nickname` VARCHAR (20) NOT NULL, `gender` INT NOT NULL, `description` TEXT NOT NULL, `email` VARCHAR(50) NOT NULL, `pwd` VARCHAR(512) NOT NULL, `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP );")
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
