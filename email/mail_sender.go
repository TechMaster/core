package email

import (
	"bytes"

	"github.com/TechMaster/core/template"
	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12/view"
)

type MailSender interface {
	SendPlainEmail(to []string, subject string, body string) error

	//Cũ SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error

	//Mới
	SendHTMLEmail(to []string, subject string, data map[string]interface{}, templateId string) error
}

type SMTPConfig struct {
	Host     string
	From     string
	Password string
	Port     int
}

var Emailer MailSender //Đối tượng phụ trách gửi email

var defaultEmailLayout string //Layout mặc định cho email

func SetDefaultEmailLayout(defaultLayout string) {
	defaultEmailLayout = defaultLayout
}

/*
Sinh HTML cho email body
data: dữ liệu
templ_layout[0]: template file
templ_layout[1]: layout file file. Nếu không có tham số này thì có 2 khả năng
A> defaultEmailLayout được cấu hình thì dùng defaultEmailLayout
B> defaultEmailLayout không được cấu thì dùng view.NoLayout
*/
func renderHTML(data map[string]interface{}, tmpl_layout ...string) (string, error) {
	var emailLayout string

	switch len(tmpl_layout) {
	case 0:
		return "", eris.New("Truyền thiếu template để sinh HTML email body")
	case 1:
		if defaultEmailLayout != "" {
			emailLayout = defaultEmailLayout
		} else {
			emailLayout = view.NoLayout
		}

	default: //2: có đầy đủ template và layout
		emailLayout = tmpl_layout[1]
	}

	viewEngine := template.ViewEngine
	buf := new(bytes.Buffer)

	if err := viewEngine.ExecuteWriter(buf, tmpl_layout[0], emailLayout, data); err != nil {
		return "", eris.NewFromMsg(err, "Lỗi generate HTML email body")
	}

	return buf.String(), nil
}
