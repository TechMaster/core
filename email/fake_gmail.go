package email

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/TechMaster/core/template"
	"github.com/TechMaster/eris"
)

/*
Dùng để kiểm thử gửi email. Email cần luôn được gửi về một địa chỉ có thể kiểm tra được receive_email
*/
type FakeGmail struct {
	config *SMTPConfig
}

var receive_email []string //Trường này sẽ đè lên to []string để luôn gửi email đến một địa chỉ cố định

func InitFakeGmail(config *SMTPConfig, test_receive_email string) {
	Emailer = FakeGmail{
		config: config,
	}
	receive_email = []string{test_receive_email}
}

//--- Hai phương thức implement interface MailSender
func (gmail FakeGmail) SendPlainEmail(to []string, subject string, body string) error {

	emailAuth := smtp.PlainAuth("me", gmail.config.From, gmail.config.Password, gmail.config.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + body)
	addr := fmt.Sprintf("%s:%d", gmail.config.Host, gmail.config.Port)

	if err := smtp.SendMail(addr, emailAuth, gmail.config.From, receive_email, msg); err != nil {
		return eris.NewFrom(err)
	}
	return nil
}

func (gmail FakeGmail) SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error {
	emailAuth := smtp.PlainAuth("me", gmail.config.From, gmail.config.Password, gmail.config.Host)

	viewEngine := template.ViewEngine
	buf := new(bytes.Buffer)
	if err := viewEngine.ExecuteWriter(buf, tmplFile, "", data); err != nil {
		return eris.NewFromMsg(err, "Lỗi generate mail body")
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + buf.String())
	addr := fmt.Sprintf("%s:%d", gmail.config.Host, gmail.config.Port)

	if err := smtp.SendMail(addr, emailAuth, gmail.config.From, receive_email, msg); err != nil {
		return eris.NewFrom(err)
	}
	return nil
}
