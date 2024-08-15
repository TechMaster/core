package rbac

import (
	"github.com/TechMaster/core/pmodel"
)

const (
	ALLOW            = "allow"
	FORBID           = "forbid"
	ALLOW_ALL        = "allow_all"
	FORBID_ALL       = "forbid_all"
	ALLOW_ONLY_ADMIN = "allow_only_admin"
)

// Dùng để in role kiểu int ra string cho dễ hiếu
var roleName = map[int]string{}

// Lưu danh sách các role
var Roles map[string]int = map[string]int{}

/*
Biểu thức hàm sẽ trả về
  - bool: true nếu là allow, false: nếu là forbid
  - danh sách role kiểu map[int]bool.
    Nếu allow thì giá trị map[int]bool đều là true
    Nếu forbid thì giá trị map[int]bool đều là false
*/
type RoleExp func() (pmodel.Roles, string)

/*
Ứng với một route = HTTP Verb + Path chúng ta có một map các role
Dùng để kiểm tra phân quyền
*/
var routesRoles = make(map[string]Route)

/*
pathsRoles có key là Path (không kèm HTTP Verb)
Dùng để in ra báo cáo cho dễ nhìn, vì các route chung một path sẽ được gom lại
*/
var pathsRoles = make(map[string]Route)

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

// Lưu cấu hình cho package RBAC
var config Config

// Cấu trúc dùng để lưu thông tin của một route
type Route struct {
	Path         string
	Method       string
	IsPrivate    bool
	Roles        pmodel.Roles
	SpecialRoles pmodel.Roles
	AccessType   string
}
