package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/TechMaster/core/template"
	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12"
)

type GmailSTMP struct {
	config *SMTPConfig
}

func InitGmail(config *SMTPConfig) {
	Emailer = GmailSTMP{
		config: config,
	}
}

//--- Hai phương thức implement interface MailSender
func (gmail GmailSTMP) SendPlainEmail(to []string, subject string, body string) error {

	emailAuth := smtp.PlainAuth("me", gmail.config.From, gmail.config.Password, gmail.config.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + body)
	addr := fmt.Sprintf("%s:%d", gmail.config.Host, gmail.config.Port)

	if err := smtp.SendMail(addr, emailAuth, gmail.config.From, to, msg); err != nil {
		return eris.NewFrom(err)
	}
	return nil
}

func (gmail GmailSTMP) SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error {
	emailAuth := smtp.PlainAuth("me", gmail.config.From, gmail.config.Password, gmail.config.Host)

	viewEngine := template.ViewEngine
	buf := new(bytes.Buffer)
	if err := viewEngine.ExecuteWriter(buf, tmplFile, iris.NoLayout, data); err != nil {
		return eris.NewFromMsg(err, "Lỗi generate mail body")
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + buf.String())
	addr := fmt.Sprintf("%s:%d", gmail.config.Host, gmail.config.Port)

	if err := smtp.SendMail(addr, emailAuth, gmail.config.From, to, msg); err != nil {
		return eris.NewFrom(err)
	}
	return nil
}
