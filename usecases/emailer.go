package usecases

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailerUseCase struct{}

func NewEmailerUseCase() *EmailerUseCase {
	return &EmailerUseCase{}
}

func (em *EmailerUseCase) SendEmailVerification(email string) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	s := fmt.Sprintf("%06d", rand.Int63n(1e6))

	emailBody := "Hello <b>You used this email to sign up for Tutorme.</i>! Please use the following code to finish signing up - {VERIFY_CODE}"

	emailBody = strings.Replace(emailBody, "{VERIFY_CODE}", s, 1)
	m := gomail.NewMessage()
	m.SetHeader("From", "arun.ranganathan111@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"arun.ranganathan111@gmail.com",
		"xvkzkvlwwdawigov",
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return s, nil
}
