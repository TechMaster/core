// Public models dùng chung giữa các dự án và module
package pmodel

import "time"

type User struct {
	tableName           struct{}  `pg:"auth.users"`
	Id                  string    `pg:",pk" json:"id"`
	Email               string    `pg:",unique" json:"email"`
	FullName            string    `valid:"required~Họ tên không được để trống,runelength(4|100)~Họ tên không hợp lệ (từ 4 - 100 ký tự)" json:"full_name"`
	Password            string    `json:"password"`
	Phone               string    `json:"phone" valid:"numeric,runelength(10|11)~Số điện thoại không hợp lệ (từ 10 - 11 ký tự)"`
	Avatar              string    `json:"avatar"`
	LinkCv              string    `json:"link_cv"`
	Description         string    `json:"description"`
	BankName            string    `json:"bank_name"`
	BankAccount         string    `json:"bank_account"`
	Slug                string    `json:"slug"`
	Roles               []int     `pg:",array" json:"roles"`
	AccessFailedCount   int32     `sql:"default:0" json:"access_failed_count"`
	EmailConfirmed      bool      `json:"email_confirmed" sql:"default:false"`
	VerifyEmailToken    string    `json:"verify_email_token"`
	VerifyEmailTokenEnd time.Time `json:"verify_email_token_end"`
	LockoutEnd          time.Time `json:"lockout_end"`
	CreatedAt           time.Time `sql:"default:now()" json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	ModifiedAt          time.Time `json:"modified_at"`
	ModifiedBy          string    `json:"modified_by"`
	UserStatus          bool      `sql:"default:true" json:"user_status"`
	Salt                string    `json:"salt"`
	NewEmail            string    `json:"new_email"`
	Dob                 time.Time `json:"dob"`
}

/*
Lưu thông tin về người đăng nhập sau khi đăng nhập thành công.
Giống hệt User nhưng loại bỏ trường HashPass
Cấu trúc này sẽ lưu vào session
*/
type AuthenInfo struct {
	UserId       string //unique id của user
	UserFullName string //họ và tên đầy đủ của user
	UserEmail    string //email cũng phải unique
	UserAvatar   string //unique id hoặc tên file ảnh đại diện
	UserPhone    string
	Roles        Roles //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
