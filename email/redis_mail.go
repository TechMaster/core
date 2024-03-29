package email

import (
	"bytes"
	"fmt"

	"html/template"

	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/eris"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

const (
	SEND_EMAIL           = "email:send"
	SEND_EMAIL_MARKETING = "email_marketing:send"
)

type EmailPayload struct {
	To       []string
	Subject  string
	Data     map[string]interface{}
	Template string
}

type MailMarketing struct {
	Sender    string
	Subject   string
	Receivers []ReceiverEmail
}

type ReceiverEmail struct {
	Email   string
	Content string
}

type RedisMail struct {
}

var Redis_mail RedisMail

var asynqClient *asynq.Client

func InitRedisMail() *asynq.Client {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Network:  viper.GetString("redis.network"),
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       1, //Do not use 0 because
	})

	Emailer = RedisMail{}
	return asynqClient
}

/*
Implement MailSender interface
*/
func (rmail RedisMail) SendPlainEmail(to []string, subject string, body string) error {

	payload, err := json.Marshal(EmailPayload{
		To:      to,
		Subject: subject,
	})
	if err != nil {
		return eris.NewFrom(err).InternalServerError()
	}

	info, err := asynqClient.Enqueue(asynq.NewTask(SEND_EMAIL, payload))
	if err != nil {
		return eris.NewFromMsg(err, "Could not enqueue task").InternalServerError()
	}
	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}

func (rmail RedisMail) SendHTMLEmail(to []string, subject string, data map[string]interface{}, templateId string) error {

	payload, err := json.Marshal(EmailPayload{
		To:       to,
		Subject:  subject,
		Data:     data,
		Template: templateId,
	})
	if err != nil {
		return eris.NewFrom(err).InternalServerError()
	}

	info, err := asynqClient.Enqueue(asynq.NewTask(SEND_EMAIL, payload))
	if err != nil {
		return eris.NewFromMsg(err, "Could not enqueue task").InternalServerError()
	}

	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}

func (rmail RedisMail) SendHTMLEmailMarketing(from, sender_name, subject string,
	to []pmodel.AuthenInfo, html_content string) (err error) {
	var emails_payload = MailMarketing{
		Sender:  sender_name + " <" + from + ">",
		Subject: subject,
	}

	var receiversEmail = []ReceiverEmail{}
	tmpl, err := template.New("name").Parse(html_content)
	if err != nil {
		logger.Log2(eris.NewFrom(err).SetType(eris.SYSERROR))
	}
	var buf *bytes.Buffer
	for _, value := range to {
		buf = bytes.NewBufferString("")
		err = tmpl.Execute(buf, map[string]interface{}{
			"Name": value.UserFullName,
		})
		if err != nil {
			logger.Log2(eris.NewFrom(err).SetType(eris.SYSERROR))
		}
		receiversEmail = append(receiversEmail, ReceiverEmail{
			Content: buf.String(),
			Email:   value.UserEmail,
		})
	}
	emails_payload.Receivers = receiversEmail

	payload, err := json.Marshal(emails_payload)
	if err != nil {
		return eris.NewFrom(err).InternalServerError()
	}

	info, err := asynqClient.Enqueue(asynq.NewTask(SEND_EMAIL_MARKETING, payload))
	if err != nil {
		return eris.NewFromMsg(err, "Could not enqueue task").InternalServerError()
	}
	fmt.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}
