package rbac

import (
	"strings"

	"github.com/TechMaster/core/pmodel"
)

/*
Cấu hình rbac
*/
func Init(configs ...Config) {
	if len(configs) == 0 {
		config = NewConfig()
	} else {
		config = configs[0]
	}
}

/*
Tạo cấu hình mặc định cho RBAC
*/
func NewConfig() Config {
	return Config{
		ForbidOverAllow:           true,
		AuthService:               "",
		MakeUnassignedRoutePublic: false,
	}
}

/*
Load các roles vào bộ nhớ

	@param fLoad: hàm trả về danh sách roles
	- Khi load các roles Phải có role admin
*/
func LoadRoles(fLoad func() []pmodel.Role) {
	roles := fLoad()
	isAdmin := false
	for _, role := range roles {
		name := strings.ToLower(role.Name)
		if name == "admin" {
			isAdmin = true
		}
		Roles[name] = role.ID
		roleName[role.ID] = name
	}
	if !isAdmin {
		panic("Phải có role admin")
	}
}

// Hàm này sẽ được gọi sau khi đã đăng ký tất cả các route

func LoadRules(fLoad func() []pmodel.Rule) {
	ConvertRules(fLoad())
}
