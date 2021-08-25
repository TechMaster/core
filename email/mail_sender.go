package email

type MailSender interface {
	SendPlainEmail(to []string, subject string, body string) error
	SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error
}

type SMTPConfig struct {
	Host     string
	From     string
	Password string
	Port     int
}

var Emailer MailSender //Đối tượng phụ trách gửi email
