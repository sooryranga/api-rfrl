package usecases

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"gopkg.in/gomail.v2"
)

type EmailerUseCase struct{}

func NewEmailerUseCase() *EmailerUseCase {
	return &EmailerUseCase{}
}

func (em *EmailerUseCase) SendEmailVerification(email string, emailType string) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	s := fmt.Sprintf("%06d", rand.Int63n(1e6))

	t := template.New("verify-html.html")

	_, b, _, _ := runtime.Caller(0)
	appDir := filepath.Dir(path.Join(path.Dir(b)))
	htmlPath := path.Join(appDir, "/assets/verify-email.html")

	htmlText, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return "", err
	}

	t, err = t.Parse(string(htmlText))

	if err != nil {
		return s, err
	}

	var tpl bytes.Buffer

	attributes, ok := tutorme.EmailTypeToEmailAttributes[emailType]

	if !ok {
		panic("email type is not found")
	}

	err = t.Execute(
		&tpl, struct {
			VerifyCode  string
			Heading     string
			Description string
		}{
			VerifyCode:  s,
			Heading:     attributes.Heading,
			Description: attributes.Description,
		})

	if err != nil {
		return s, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "arun.ranganathan111@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"arun.ranganathan111@gmail.com",
		"xvkzkvlwwdawigov",
	)

	if err := d.DialAndSend(m); err != nil {
		return s, err
	}

	return s, nil
}
