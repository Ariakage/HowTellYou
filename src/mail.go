package main

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

func SendEmail(senderAddr string, receviceAddr string, server string, port int, authcode string, content string) error {
	auth := smtp.PlainAuth("", senderAddr, authcode, server)
	to := []string{receviceAddr}
	nickname := ""
	subject := ""
	contentType := "Content-Type:text/html;charset=UTF-8\r\n"
	// 支持群发
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + senderAddr + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + content)
	err := smtp.SendMail(server+":"+strconv.Itoa(port), auth, senderAddr, to, msg)
	if err != nil {
		err = fmt.Errorf("Send Err%v", err)
		fmt.Println(err)
		return err
	}
	return nil
}
