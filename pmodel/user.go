// Public models dùng chung giữa các dự án và module
package pmodel

import "time"

type User struct {
	tableName           struct{}  `pg:"auth.users"` //Postgresql không chấp nhận bảng có tên là user
	Id                  string    `pg:",pk" json:"id"`                //chuỗi ngẫu nhiên duy nhất
	Email               string    `pg:",unique" json:"email"`
	FullName            string    `valid:"required~Họ tên không được để trống,runelength(4|100)~Họ tên không hợp lệ (từ 4 - 100 ký tự)" json:"full_name"`
	Password            string    `json:"password"`     // Hashed password. Tuyệt đối không lưu raw password
	Phone               string    `json:"phone" valid:"numeric,runelength(10|11)~Số điện thoại không hợp lệ (từ 10 - 11 ký tự)" json:"full_name"`        // Số di động ở VN có từ 10-11 chữ số
	Avatar              string    `json:"avatar"`       // Ảnh đại diện
	Description         string    `json:"description"`  // Mô tả
	BankName            string    `json:"bank_name"`    // Tên ngân hàng
	BankAccount         string    `json:"bank_account"` // Số tài khoản ngân hàng
	Slug                string    `json:"slug"`
	Roles               []int     `pg:",array" json:"roles"`
	AccessFailedCount   int32     `sql:"default:0" json:"access_failed_count"` // Số lần đăng nhập sai, mặc định là 0
	EmailConfirmed      bool      `json:"email_confirmed" sql:"default:false"` // Email đã xác nhận (kích hoạt) hay chưa
	VerifyEmailToken    string    `json:"verify_email_token"`                  // Token để xác thực Email
	VerifyEmailTokenEnd time.Time `json:"verify_email_token_end"`              // Thời gian hiệu lực của Token xác thực email
	LockoutEnd          time.Time `json:"lockout_end"`                         // Thời điểm hết khoá tài khoản
	CreatedAt           time.Time `sql:"default:now()" json:"created_at"`      // Ngày tài khoản được tạo
	CreatedBy           string    `json:"created_by"`                          // Id người tạo tài khoản, Null là người dùng tự đăng ký tài khoản
	ModifiedAt          time.Time `sql:"default:now()" json:"modified_at"`     // Ngày gần nhất tài khoản cập nhật thông tin
	ModifiedBy          string    `json:"modified_by"`                         // Người cập nhật thông tin tài khoản gần nhất
	UserStatus          bool      `sql:"default:true" json:"user_status"`      // True là active, False là unactive, mặc định là True
	Salt                string    `json:"salt"`                                // Dùng để kiểm tra hash password. Khi dùng BCrypt không cần nữa
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
	UserAvatar  string //unique id hoặc tên file ảnh đại diện
	UserPhone	string
	Roles        Roles  //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
