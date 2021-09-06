package email

import (
	"bytes"
	"time"

	"github.com/TechMaster/core/template"
	"github.com/kataras/iris/v12"

	"github.com/TechMaster/eris"
	"github.com/go-pg/pg/v10"
)

type EmailDB struct {
	db *pg.DB
}

func InitEmailDB(db_ *pg.DB) {
	Emailer = EmailDB{
		db: db_,
	}

}

type EmailStore struct {
	tableName  struct{} `pg:"debug.emailstore"`
	Id         int      `pg:",pk"`
	Receipient string
	Subject    string
	Body       string
	CreatedAt  time.Time
}

func (emailDB EmailDB) SendPlainEmail(to []string, subject string, body string) error {
	emailitem := EmailStore{
		Receipient: to[0],
		Subject:    subject,
		Body:       body,
	}
	if _, err := emailDB.db.Model(&emailitem).Insert(); err != nil {
		return eris.NewFrom(err).InternalServerError()
	}
	return nil
}

func (emailDB EmailDB) SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error {
	viewEngine := template.ViewEngine
	buf := new(bytes.Buffer)
	if err := viewEngine.ExecuteWriter(buf, tmplFile, iris.NoLayout, data); err != nil {
		return eris.NewFromMsg(err, "Lỗi generate mail body")
	}

	emailitem := EmailStore{
		Receipient: to[0],
		Subject:    subject,
		Body:       buf.String(),
	}
	if _, err := emailDB.db.Model(&emailitem).Insert(); err != nil {
		return eris.NewFrom(err).InternalServerError()
	}
	return nil
}

func (emailDB EmailDB) GetMail(to string) (email *EmailStore, err error) {
	email = new(EmailStore)
	if _, err := emailDB.db.Query(email, `SELECT * FROM debug.emailstore 
	WHERE receipient = ?
	ORDER BY id DESC LIMIT 1`, to); err != nil {
		return nil, eris.NewFrom(err).InternalServerError()
	}
	return email, nil
}
