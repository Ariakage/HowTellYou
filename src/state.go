package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

func GenerateLoginToken(chars string, salt_length int, key string, user_id int, login_time time.Time) (string, string) {
	salt := GenerateRandomString(salt_length, chars)
	t1 := strconv.Itoa(user_id) + "@" + strconv.FormatInt(login_time.UnixNano(), 10) + "|" + salt
	t2 := base64.StdEncoding.EncodeToString([]byte(t1))
	res, err := SCEncryptString(t2, key, "aes")
	if err != nil {
		fmt.Println(err)
	}
	return res, salt
}
