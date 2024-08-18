package rbac

import (
	"fmt"
	"strings"

	"github.com/TechMaster/core/pmodel"

	"github.com/kataras/iris/v12"
)

/*
Lập ra danh sách các public route. Đây là những route không được đăng ký qua rbac
Hãy đặt hàm này ở cuối cùng hàm main. Ngay trước câu lệnh listen port để đảm bảo
nó quét được hết tất cả các public route
*/

func BuildPublicRoute(app *iris.Application) {
	for _, route := range app.GetRoutes() {
		paddingRoute := correctRoute(route.Name)

		//Nếu IsPrivate = false thì đây là route public
		if r, ok := routesRoles[paddingRoute]; ok && !r.IsPrivate {
			publicRoutes[paddingRoute] = true
		}
	}
}

/*
In ra danh sách những đường dẫn public không kiểm tra quyền
*/
func DebugPublicRouteRole() {
	fmt.Println("*** Public Routes ***")
	fmt.Println("Total", "----", len(publicRoutes))
	for route := range publicRoutes {
		fmt.Println(route)
	}
}

/*
In ra thông tin Private Route - Role
route = HTTP Verb + path
*/
func DebugRouteRole() {
	fmt.Println("*** Private Routes ***")
	for path, route := range routesRoles {
		fmt.Println("-"+path, "-", route.IsPrivate)
		for role, allow := range route.Roles {
			if allow.(bool) {
				fmt.Println("      " + RoleName(role) + "allow") //allow role in trắng
			} else {
				fmt.Println("     \033[31m^" + RoleName(role) + "\033[0m" + "forbid") //forbid role in đỏ
			}
		}
	}
}

/*
In ra thông tin debug Route - Rote thành 2 phần
Phần Public: những route không cần kiểm tra phân quyền
Phần Private: những route cần kiểm tra quyền
*/
func DebugPathRole() {
	fmt.Println("*** Private Path ***")
	for path, route := range pathsRoles {
		fmt.Println("-" + path + "-" + route.Method)
		for role := range route.Roles {
			fmt.Println("     " + RoleName(role))
		}
	}
}

/*
	Insert space between HTTP Verb and Path

Input: GET/blog
Output: GET /blog
*/
func correctRoute(route string) string {
	posFirstSlash := strings.Index(route, "/")
	return route[0:posFirstSlash] + " " + route[posFirstSlash:]
}

// Chuyển role từ in string
func RoleName(role int) string {
	return roleName[role]
}

/*
Chuyển roles kiểu map[int]bool thành mảng string mô tả các role`
*/
func RolesNames(roles pmodel.Roles) (rolesNames []string) {
	for role := range roles {
		rolesNames = append(rolesNames, RoleName(role))
	}
	return
}

// Hàm chuyển đổi từ Rule sang Route
func ConvertRules(rules []pmodel.Rule) {

	for _, rule := range rules {
		route := Route{
			IsPrivate:  rule.IsPrivate,
			Path:       rule.Path,
			Method:     strings.ToUpper(rule.Method),
			AccessType: rule.AccessType,
		}
		switch strings.ToLower(rule.AccessType) {
		case ALLOW:
			roles, _ := Allow(rule.Roles...)()
			route.Roles = roles
		case ALLOW_ALL:
			roles, _ := AllowAll()()
			route.Roles = roles
		case FORBID:
			roles, _ := Forbid(rule.Roles...)()
			route.Roles = roles
		case FORBID_ALL:
			roles, _ := ForbidAll()()
			route.Roles = roles
		case ALLOW_ONLY_ADMIN:
		default:
			roles, _ := AllowOnlyAdmin()()
			route.Roles = roles
		}
		routesRoles[route.Method+" "+route.Path] = route
	}
}

/*
	Load lại các rules public, dùng khi có thay đổi về rules từ database
*/

func ReloadPublicRoute() {
	for path, route := range routesRoles {
		if _, ok := publicRoutes[path]; ok && route.IsPrivate {
			delete(publicRoutes, path)
		}
		if !route.IsPrivate {
			publicRoutes[path] = true
		}
	}
}

/*
	Hàm tự động insert các rules vào database
	@param funcInsert: hàm insert rules vào database
*/

func AutoRegisterRules(funcInsert func(rules []pmodel.Rule)) {
	rules := []pmodel.Rule{}
	for _, route := range routesRoles {
		rule := pmodel.Rule{
			Path:       route.Path,
			Method:     route.Method,
			AccessType: route.AccessType,
			IsPrivate:  route.IsPrivate,
			Service:    config.Service,
		}
		rules = append(rules, rule)
	}
	funcInsert(rules)
}

func AutoRegisterRoles(funcInsert func(roles []pmodel.Role), rolesDefault []string) {
	roles := []pmodel.Role{}
	for _, role := range rolesDefault {
		roles = append(roles, pmodel.Role{
			Name: strings.ToLower(role),
		})
	}
	funcInsert(roles)
}
