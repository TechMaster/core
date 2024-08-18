package rbac

import (
	"github.com/TechMaster/core/pmodel"
)

// Danh sách các role có thể truy xuất
func Allow(roles ...int) RoleExp {
	return func() (pmodel.Roles, string) {
		mapRoles := make(pmodel.Roles)
		for _, role := range roles {
			mapRoles[role] = true
		}
		return mapRoles, ALLOW
	}
}

// Cho phép tất cả các role
func AllowAll() RoleExp {
	return func() (pmodel.Roles, string) {
		mapRoles := make(pmodel.Roles)
		for _, role := range Roles {
			mapRoles[role] = true
		}
		return mapRoles, ALLOW_ALL
	}
}


// Danh sách các role bị cấm truy cập
func Forbid(roles ...int) RoleExp {
	return func() (pmodel.Roles, string) {
		mapRoles := make(pmodel.Roles)
		for _, role := range roles {
			mapRoles[role] = false
		}
		return mapRoles, FORBID
	}
}

// Cấm tất cả các role ngoại trừ root
func ForbidAll() RoleExp {
	return func() (pmodel.Roles, string) {
		mapRoles := make(pmodel.Roles)
		for _, role := range Roles {
			mapRoles[role] = false
		}
		return mapRoles, FORBID_ALL
	}
}
