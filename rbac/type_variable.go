package rbac

import "github.com/TechMaster/core/pmodel"

/*
Nếu thêm sửa xoá role thì cập nhật danh sách const này
*/
const (
	ROOT       = 0 //Role đặc biệt, vượt qua mọi logic kiểm tra quyền khi config.RootAllow = true
	ADMIN      = 1
	STUDENT    = 2
	TRAINER    = 3
	SALE       = 4
	EMPLOYER   = 5
	AUTHOR     = 6
	EDITOR     = 7
	MAINTAINER = 8
	SYSOP      = 9 //Role mới
)

//Mảng này phải tương ứng với danh sách const khai báo ở trên
var allRoles = []int{ROOT, ADMIN, STUDENT, TRAINER, SALE, EMPLOYER, AUTHOR, EDITOR, MAINTAINER, SYSOP}

//Dùng để in role kiểu int ra string cho dễ hiếu
var roleName = map[int]string{
	ROOT:       "root",
	ADMIN:      "admin",
	STUDENT:    "student",
	TRAINER:    "trainer",
	SALE:       "sale",
	EMPLOYER:   "employer",
	AUTHOR:     "author",
	EDITOR:     "editor",
	MAINTAINER: "maintainer",
	SYSOP:      "system operator",
}

/*
Biểu thức hàm sẽ trả về
- bool: true nếu là allow, false: nếu là forbid
- danh sách role kiểu map[int]bool.
  Nếu allow thì giá trị map[int]bool đều là true
	Nếu forbid thì giá trị map[int]bool đều là false
*/
type RoleExp func() pmodel.Roles

/*
Ứng với một route = HTTP Verb + Path chúng ta có một map các role
Dùng để kiểm tra phân quyền
*/
var routesRoles = make(map[string]pmodel.Roles)

/*
pathsRoles có key là Path (không kèm HTTP Verb)
Dùng để in ra báo cáo cho dễ nhìn, vì các route chung một path sẽ được gom lại
*/
var pathsRoles = make(map[string]HTTPVerbRoles)

/*
kiểu HTTPVerbRoles là map có key là 'GET', 'POST', 'PUT', 'DELETE'
Value là map các role
HTTPVerbRoles dùng để gom các roles gán cho từng HTTP Verb ứng với một path
*/
type HTTPVerbRoles map[string]pmodel.Roles

/*
Danh sách các public routes dùng trong hàm CheckPermission
*/
var publicRoutes = make(map[string]bool)

/*
Cấu hình cho hệ thống RBAC
*/
type Config struct {
	/* Nếu một người có 2 role A và B. Ở route X, role A bị cấm và role B được phép.
	ForbidOverAllow = true (mặc định) thì người đó bị cấm ở route X
	ForbidOverAllow = false thì người đó được phép ở route X
	*/
	ForbidOverAllow bool

	/* Nếu RootAllow bằng true thì người có role là Root sẽ
	vượt qua tất cả logic kiểm tra permission.
	Mặc định RootAllow = false
	*/
	RootAllow bool

	/*
		Đường dẫn đến AuthService trong mạng Docker Compose, Docker Swarm
		"" nếu RBAC không cần kết nối đến AuthService
	*/
	AuthService string

	/* MakeUnassignedRoutePublic = true sẽ biến tất cả những đường dẫn
	không có trong map routesRoles mặc nhiên là public
	mặc định là false
	*/
	MakeUnassignedRoutePublic bool
}

//Lưu cấu hình cho package RBAC
var config Config
