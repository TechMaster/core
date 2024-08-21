package pmodel

/*
	Danh sách các role gán cho một người hoặc một route

dữ liệu trong value kiểu bool nhưng tôi không dùng bool mà dùng interface{}
bởi nếu dùng map[int]bool khi truyền vào key không tồn tại luôn trả về false
nhưng tôi mong muốn phải trả về nil mới đúng bản chất
*/
type Roles map[int]interface{}

/*
Chuyển đổi Roles kiểu map[int] bool sang mảng []int để lưu xuống CSDL
*/
func RolesToIntArr(roles Roles) []int {
	keys := make([]int, 0, len(roles))
	for k := range roles {
		keys = append(keys, k)
	}
	return keys
}

/*
Chuyển đổi kiểu intArray trong đó mỗi phần tử ứng với một role, sang kiểu map[int] bool
*/
func IntArrToRoles(intArr []int) Roles {
	roles := make(Roles)
	for _, role := range intArr {
		roles[role] = true
	}
	return roles
}

/*
Role là cấu hình cho việc kiểm tra quyền truy cập
*/
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // Tên của role | lowercase
}
