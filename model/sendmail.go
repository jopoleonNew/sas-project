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
	//linkTag := "<a href='" + link + "'>" + link + "</a>"
	//body := "<a target='blank' href='" + link + "'>" + link + "/>"
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
	//linkTag := "<a href='" + link + "'>" + link + "</a>"

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

//func SendEmailwithKey(username, key, emailaddr, requestHost string) error {
//	//fmt.Println("Start time: ", time.Now())
//	//Sendkey2017
//
//	emailUser := &EmailUser{"sasmailserver", "SendSASKey2017", "smtp.gmail.com", 587}
//
//	auth := smtp.PlainAuth("",
//		emailUser.Username,
//		emailUser.Password,
//		emailUser.EmailServer,
//	)
//
//	var err error
//
//	link := requestHost + "/activateuser?username=" + username + "&activationkey=" + key
//	//linkTag := "<a href='" + link + "'>" + link + "</a>"
//	body := "<a target=«blank» href=«" + link + "»>" + link + "/>"
//	msg := []byte("To: " + emailaddr + "\r\n" +
//		"Subject: Activation letter from SASmailServer\r\n" +
//		"\r\n" +
//		"This is your activation link: \r\n\n" + body)
//	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
//		auth,
//		emailUser.Username,
//		[]string{emailaddr},
//		msg)
//	if err != nil {
//		log.Print("model.SendEmailwithKey error: ", err)
//		return err
//	}
//	return nil
//	//fmt.Println("End time: ", time.Now())
//}
// type RoomInfo struct {
// 	Username string `json:Username bson:Username`
// 	Location string `json:Location bson:Location`
// 	Title    string `json:Title bson:Title`
// 	Cost     string `json:Cost bson:Cost`
// 	ImageUrl string `json:ImageUrl bson:ImageUrl`
// }
//
//func SendEmailTrackList(roominfoslice []models.RoomInfo, emailadr, user, loca string) {
//
//	type EmailUser struct {
//		Username    string
//		Password    string
//		EmailServer string
//		Port        int
//	}
//
//	emailUser := &EmailUser{"spareroommailserver", "Sendkey2017", "smtp.gmail.com", 587}
//
//	auth := smtp.PlainAuth("",
//		emailUser.Username,
//		emailUser.Password,
//		emailUser.EmailServer,
//	)
//
//	var err error
//	var msgbody string
//	for _, info := range roominfoslice {
//		msgbody += " " + info.Title + " " + info.Cost + " " + "https:" + info.ImageUrl + "\n\n"
//	}
//
//	msg := []byte("To: " + emailadr + "\r\n" +
//		"Subject: New rooms in your tracking area: " + loca + "\r\n" +
//		"\r\n" + " Hello " + user + "\n Here are some new rooms for you: \n \n" +
//		msgbody)
//	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
//		auth,
//		emailUser.Username,
//		[]string{emailadr},
//		msg)
//	if err != nil {
//		log.Print("ERROR: attempting to send a mail ", err)
//	}
//	//fmt.Println("End time: ", time.Now())
//}
