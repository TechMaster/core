package email

import (
	"bytes"
	"fmt"

	"github.com/TechMaster/core/template"
	"github.com/TechMaster/eris"
	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"

	"github.com/hibiken/asynq"
)

const (
	SEND_EMAIL = "email:send"
)

type EmailPayload struct {
	To      []string
	Subject string
	Msg     []byte
}

type RedisMail struct {
}

var asynqClient *asynq.Client

func InitRedisMail() *asynq.Client {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Network:  viper.GetString("redis.network"),
		Addr:     viper.GetString("redis.address"),
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

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + body)

	payload, err := json.Marshal(EmailPayload{
		To:      to,
		Subject: subjectStr,
		Msg:     msg,
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

func (rmail RedisMail) SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error {
	viewEngine := template.ViewEngine
	buf := new(bytes.Buffer)
	if err := viewEngine.ExecuteWriter(buf, tmplFile, iris.NoLayout, data); err != nil {
		return eris.NewFromMsg(err, "Lỗi generate mail body")
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectStr := "Subject: " + subject + "!\n"
	msg := []byte(subjectStr + mime + "\n" + buf.String())

	payload, err := json.Marshal(EmailPayload{
		To:      to,
		Subject: subjectStr,
		Msg:     msg,
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
