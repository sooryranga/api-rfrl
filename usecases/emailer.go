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

	"gopkg.in/gomail.v2"
)

type EmailerUseCase struct{}

func NewEmailerUseCase() *EmailerUseCase {
	return &EmailerUseCase{}
}

func (em *EmailerUseCase) SendEmailVerification(email string) (string, error) {
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

	err = t.Execute(
		&tpl, struct {
			VerifyCode string
		}{
			VerifyCode: s,
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
