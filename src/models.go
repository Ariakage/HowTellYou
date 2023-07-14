/*
CREATED BY DR.ALANORAGE on 2023.07.11
*/
package main

import "time"

type User struct {
	id          int
	name        string
	email       string
	password    string
	create_time time.Time
}
