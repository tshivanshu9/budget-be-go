package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"gopkg.in/gomail.v2"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

func NewMailer() Mailer {
	mailPort := os.Getenv("MAIL_PORT")
	mailHost := os.Getenv("MAIL_HOST")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")
	mailPortInt, err := strconv.Atoi(mailPort)

	if err != nil {
		panic("Invalid mail port")
	}

	d := gomail.NewDialer(mailHost, mailPortInt, mailUsername, mailPassword)
	return Mailer{
		dialer: d,
		sender: mailSender,
	}
}

func (m *Mailer) Send(recipient string, templateFile string, data EmailData) error {
	absolutePath := filepath.Join("templates", templateFile)
	tmpl, err := template.ParseFS(templateFS, absolutePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("To", recipient)
	gomailMessage.SetHeader("From", m.sender)
	gomailMessage.SetHeader("subject", subject.String())
	gomailMessage.SetBody("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(gomailMessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
