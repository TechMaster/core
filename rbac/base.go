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

func LoadRoles(fLoad func() []pmodel.Role, rolesRequire ...map[string]bool) {
	roles := fLoad()
	// Tạo một bản sao của rolesRequire để theo dõi các roles chưa được tìm thấy
	missingRoles := make(map[string]bool)
	if config.HighestRole != "" {
		missingRoles[config.HighestRole] = true
	} else {
		missingRoles[DEFAULT_HIGHEST_ROLE] = true
	}
	if len(rolesRequire) > 0 {
		for k, v := range rolesRequire[0] {
			missingRoles[k] = v
		}
	}
	// Duyệt qua các roles được tải
	for _, role := range roles {
		name := strings.ToLower(role.Name)
		if ok := missingRoles[name]; ok {
			delete(missingRoles, name) // Xóa role khỏi missingRoles nếu tìm thấy
		}
		Roles[name] = role.ID
		roleName[role.ID] = name
	}

	// Kiểm tra nếu có role yêu cầu nào không được tìm thấy
	if len(missingRoles) > 0 {
		var builder strings.Builder
		builder.WriteString("Phải có role: ")
		for name := range missingRoles {
			builder.WriteString(name + ", ")
		}
		// Xóa dấu phẩy cuối cùng và khoảng trắng
		str := builder.String()
		if len(str) > 0 {
			str = str[:len(str)-2]
		}
		panic(str)
	}
}

// Hàm này sẽ được gọi sau khi đã đăng ký tất cả các route

func LoadRules(fLoad func() []pmodel.Rule) {
	ConvertRules(fLoad())
}
