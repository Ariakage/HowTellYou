/*
CREATED BY DR.ALANORAGE on 2023.07.15
*/
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

func closeDB(db *sql.DB) {
	db.Close()
}
