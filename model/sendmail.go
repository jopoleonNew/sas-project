package model

import (
	"log"
	"net/smtp"
	"strconv"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

// info : SASmailServer
// adress: sasmailserver@gmail.com
// password: SendSASKey2017
func SendEmailwithKey(username, key, emailaddr, requestHost string) error {
	link := "https://" + requestHost + "/activateuser?username=" + username + "&activationkey=" + key
	msg := []byte("To: " + emailaddr + "\r\n" +
		"Subject: Activation letter from SASmailServer\r\n" +
		"\r\n" +
		"This is your activation link: \r\n" + link)
	err := SendEmail(emailaddr, msg)
	if err != nil {
		log.Println("model.SendEmailwithKey SendEmail error: ", err)
		return err
	}
	return nil
}
func SendEmailRestorePass(linkkey, emailaddr, requestHost string) error {

	link := "https://" + requestHost + "/restorepass?secretkey=" + linkkey

	msg := []byte("To: " + emailaddr + "\r\n" +
		"Subject: Restore password letter SASmailServer\r\n" +
		"\r\n" +
		"Обратите внимание, что ссылка в этом письме будет активна только 5 часов, и ссылку в письме можно использовать только один раз \r\n" +
		"Чтобы восстановить пароль, перейдите по ссылке и следуете дальнейшим указаниям: \r\n\n" + link)
	err := SendEmail(emailaddr, msg)
	if err != nil {
		log.Println("model.SendEmailwithKey SendEmail error: ", err)
		return err
	}
	return nil
}

func SendEmail(emailaddr string, msgbody []byte) error {
	log.Println("SendEmail used with: ", emailaddr, string(msgbody))
	emailUser := &EmailUser{"sasmailserver", "SendSASKey2017", "smtp.gmail.com", 587}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)
	var err error
	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		auth,
		emailUser.Username,
		[]string{emailaddr},
		msgbody)
	if err != nil {
		log.Println("modl.SendEmail() error: ", err)
		return err
	}
	return nil
}
