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
	Id       string            //unique id của user
	Password string            //password đã được băm
	FullName string            //họ và tên đầy đủ của user
	Email    string            //email cũng phải unique
	Phone    string            //số di động của user
	Avatar   string            //unique id hoặc tên file ảnh đại diện
	Roles    Roles             //kiểu map[int]bool
	Attrs    map[string]string `pg:",hstore"` //sử dụng kiểu store để lưu dữ liệu kiểu map[string]string
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
	Roles    Roles  //kiểu map[int]bool
}
