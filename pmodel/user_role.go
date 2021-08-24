// Public models dùng chung giữa các dự án và module
package pmodel

/* Danh sách các role gán cho một người hoặc một route
dữ liệu trong value kiểu bool nhưng tôi không dùng bool mà dùng interface{}
bởi nếu dùng map[int]bool khi truyền vào key không tồn tại luôn trả về false
nhưng tôi mong muốn phải trả về nil mới đúng bản chất
*/
type Roles map[int]interface{}

//Thông tin tài khoản
type User struct {
	User  string
	Pass  string
	Email string
	Roles Roles
}

/*
Lưu thông tin về người đăng nhập sau khi đăng nhập thành công.
Cấu trúc này sẽ lưu vào session
*/
type AuthenInfo struct {
	User  string
	Email string
	Roles Roles //kiểu map[int]bool
}
