package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL     string
	Subject string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func (cfg *apiConfig) SendEmail(user *database.User, data *EmailData, emailTemp string) error {

	from := cfg.EmailFrom
	smtpPass := cfg.SMTPPass
	smtpUser := cfg.SMTPUser
	to := user.Email
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return err
	}

	err = template.ExecuteTemplate(&body, emailTemp, &data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
