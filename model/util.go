package model

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateKey32chars() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err) // out of randomness, should never happen
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
