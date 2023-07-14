package main

import (
	"fmt"
	"time"
)

func generateLoginToken(key string, user_id int, login_time time.Time) (string, string) {

	res, err := SCEncryptString("", key, "aes")
	if err != nil {
		fmt.Println(err)
	}
	return res, ""
}
