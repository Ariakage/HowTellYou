package main

import "time"

type User struct {
	id          int
	name        string
	email       string
	password    string
	create_time time.Time
}
