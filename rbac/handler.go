package rbac

import (
	"net/http"
	"regexp"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
)

/*
Gán role vào route và path
route = HTTP method + path
roles là kết quả trả về từ hàm kiểu roleExp()
*/
func assignRoles(route Route) {
	//Sử dụng regex để thay thế double slashes bằng single slash
	re, _ := regexp.Compile("/+")
	route.Path = re.ReplaceAllLiteralString(route.Path, "/")

	routesRoles[route.Method+" "+route.Path] = route

	if _, ok := pathsRoles[route.Path]; !ok {
		pathsRoles[route.Path] = route
	}
}

func Get(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		Method:     http.MethodGet,
		IsPrivate:  isPrivate,
		Roles:      roles,
		AccessType: accessType,
	}
	party.Handle(route.Method, relativePath, handlers...)
	assignRoles(route)
}

func Post(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		Method:     http.MethodPost,
		IsPrivate:  isPrivate,
		Roles:      roles,
		AccessType: accessType,
	}

	party.Handle(route.Method, relativePath, handlers...)

	assignRoles(route)
}

func Put(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		Method:     http.MethodPut,
		IsPrivate:  isPrivate,
		Roles:      roles,
		AccessType: accessType,
	}
	party.Handle(route.Method, relativePath, handlers...)
	assignRoles(route)
}

func Delete(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		Method:     http.MethodDelete,
		IsPrivate:  isPrivate,
		Roles:      roles,
		AccessType: accessType,
	}
	party.Handle(route.Method, relativePath, handlers...)
	assignRoles(route)
}

func Patch(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		Method:     http.MethodPatch,
		Roles:      roles,
		AccessType: accessType,
	}
	party.Handle(route.Method, relativePath, handlers...)
	assignRoles(route)
}

func Any(party router.Party, relativePath string, roleExp RoleExp, isPrivate bool, handlers ...context.Handler) {
	roles, accessType := roleExp()
	route := Route{
		Path:       party.GetRelPath() + relativePath,
		IsPrivate:  isPrivate,
		Roles:      roles,
		AccessType: accessType,
	}
	party.Any(relativePath, handlers...)
	for _, method := range router.AllMethods {
		route.Method = method
		assignRoles(route)
	}
}
