package pmodel

/*
Rule là cấu hình cho việc kiểm tra quyền truy cập
- Nếu IsPrivate = true thì cần kiểm tra quyền và Roles, AccessType không có ý nghĩa, có thể để AccessType = "allow_all" cho dễ hiểu
*/
type Rule struct {
	ID           int    //ID của rule
	Name         string //Tên của rule
	Roles        []int  `pg:",array"` //Danh sách các role có thể truy xuất
	AccessType   string //Allow, AllowAll, AllowOnlyAdmin, Forbid, ForbidAll
	Method       string //GET, POST, PUT, DELETE, PATCH
	Path         string //Đường dẫn
	IsPrivate    bool   //true: cần kiểm tra quyền, false: không cần kiểm tra quyền
	Service     string //Dịnh nghĩa các rule cho các service khác nhau
}
