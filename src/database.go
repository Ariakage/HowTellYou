package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func openDataBase(fp string) *sql.DB {
	db, err := sql.Open("mysql", fp)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func execSQL(db *sql.DB, cmd string) sql.Result {
	res, err := db.Exec(cmd)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func closeDB(db *sql.DB) {
	db.Close()
}
