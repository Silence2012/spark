package utils

import (
	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
	"../constants"
)

func SendEmail(from string, to []string, cc string, subject string, contentType string, body string, attachments ...string ) error {

	mailSender := gomail.NewMessage()
	mailSender.SetHeader("From", from)
	mailSender.SetHeader("To", to...)

	mailSender.SetHeader("Subject", subject)
	mailSender.SetBody(contentType, body)

	for _, attachment := range attachments {
		mailSender.Attach(attachment)
	}

	emailHost := beego.AppConfig.String(constants.EmailHost)
	emailPort, convertErr := beego.AppConfig.Int(constants.EmailPort)
	if convertErr != nil {
		//默认587
		emailPort = 587
	}
	emailUser := beego.AppConfig.String(constants.EmailUser)
	emailPwd := beego.AppConfig.String(constants.EmailPwd)
	d := gomail.NewDialer(emailHost, emailPort, emailUser, emailPwd)

	if err := d.DialAndSend(mailSender); err != nil {
		return err
	}

	return nil
}

