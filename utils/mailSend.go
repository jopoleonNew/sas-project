// Package utils implements some utility fucntions like sending email, generating ID etc.

package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

func GenerateKey32chars() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return fmt.Sprintf("%x", buf)
	// or hex.EncodeToString(buf)
	// or base64.StdEncoding.EncodeToString(buf)
}

//https://www.socketloop.com/tutorials/golang-validate-email-address-with-regular-expression
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func SendEmailwithKey(key, addres, currentActivationUrl string) {
	fmt.Println("Start time: ", time.Now())
	//Sendkey2017
	type EmailUser struct {
		Username    string
		Password    string
		EmailServer string
		Port        int
	}

	emailUser := &EmailUser{"spareroommailserver", "Sendkey2017", "smtp.gmail.com", 587}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)

	var err error

	link := currentActivationUrl + "?" + key
	//linkTag := "<a href='" + link + "'>" + link + "</a>"
	msg := []byte("To: " + addres + "\r\n" +
		"Subject: Activation letter from SpareRoomScraper\r\n" +
		"\r\n" +
		"This is your activation link: \r\n" + link)
	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		auth,
		emailUser.Username,
		[]string{addres},
		msg)
	if err != nil {
		log.Print("ERROR: attempting to send a mail ", err)
	}
	fmt.Println("End time: ", time.Now())
}
