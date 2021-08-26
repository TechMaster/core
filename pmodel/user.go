// Public models dùng chung giữa các dự án và module
package pmodel

import "time"

type User struct {
	tableName           struct{} `pg:"auth.users"` //Postgresql không chấp nhận bảng có tên là user
	Id                  string   `pg:",pk"`        //chuỗi ngẫu nhiên duy nhất
	Email               string   `pg:",unique"`
	FullName            string   `valid:"required~Họ tên không được để trống,runelength(4|100)~Họ tên không hợp lệ (từ 4 - 100 ký tự)"`
	Password            string   // Hashed password. Tuyệt đối không lưu raw password
	Phone               string   // Số di động ở VN có từ 10-11 chữ số
	Avatar              string   // Ảnh đại diện
	Description         string   // Mô tả
	BankName            string   // Tên ngân hàng
	BankAccount         string   // Số tài khoản ngân hàng
	Slug                string
	Roles               []int     `pg:",array"`
	AccessFailedCount   int32     `sql:"default:0"`                            // Số lần đăng nhập sai, mặc định là 0
	EmailConfirmed      bool      `json:"email_confirmed" sql:"default:false"` // Email đã xác nhận (kích hoạt) hay chưa
	VerifyEmailToken    string    `json:"verify_email_token"`                  // Token để xác thực Email
	VerifyEmailTokenEnd time.Time // Thời gian hiệu lực của Token xác thực email
	LockoutEnd          time.Time // Thời điểm hết khoá tài khoản
	CreatedAt           time.Time `sql:"default:now()"` // Ngày tài khoản được tạo
	CreatedBy           string    // Id người tạo tài khoản, Null là người dùng tự đăng ký tài khoản
	ModifiedAt          time.Time `sql:"default:now()"` // Ngày gần nhất tài khoản cập nhật thông tin
	ModifiedBy          string    // Người cập nhật thông tin tài khoản gần nhất
	UserStatus          bool      `sql:"default:true"` // True là active, False là unactive, mặc định là True
	Salt                string    // Dùng để kiểm tra hash password. Khi dùng BCrypt không cần nữa
}

/*
Lưu thông tin về người đăng nhập sau khi đăng nhập thành công.
Giống hệt User nhưng loại bỏ trường HashPass
Cấu trúc này sẽ lưu vào session
*/
type AuthenInfo struct {
	Id       string //unique id của user
	FullName string //họ và tên đầy đủ của user
	Email    string //email cũng phải unique
	Avatar   string //unique id hoặc tên file ảnh đại diện
	Roles    Roles  //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
