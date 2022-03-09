package email

import (
	"testing"

	"github.com/TechMaster/core/config"
	"github.com/TechMaster/core/template"
	"github.com/kataras/iris/v12"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

/*
Kiểm thử  chức năng  gửi  email qua Redis Stream và AsynQ
github.com/hibiken/asynq
Hãy chạy ở chế độ debug test
*/
func Test_Send_Email_Send(t *testing.T) {
	config.ReadConfig("..")
	asynClient := InitRedisMail()
	defer asynClient.Close()
	var err error

	for i := 0; i < 3; i++ {
		err = Emailer.SendPlainEmail([]string{"cuong@techmaster.vn"}, gofakeit.Sentence(10), gofakeit.Paragraph(3, 5, 7, "\n\n"))
		if err != nil {
			break
		}
	}
	assert := assert.New(t)
	assert.Nil(err)
}

func Test_Send_Email_Marketing(t *testing.T) {
	config.ReadConfig("..")
	asynClient := InitRedisMail()
	defer asynClient.Close()
	app := iris.New()
	template.InitBlockEngine(app, "../views", "default")

	var redis_mail RedisMail
	err := redis_mail.SendHTMLEmailMarketing("ba@techmaster.vn", "Xuân Ba", "Đức test", "",
	[]string{"nhatduc@techmaster.vn", "nhatduc.hoanghapaper@gmail.com"}, map[string]interface{}{
	}, "index")
	assert := assert.New(t)
	assert.Nil(err)
}
