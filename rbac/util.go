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
		//Nếu không tìm thấy padddingRoute trong map những private route
		if roles := routesRoles[paddingRoute]; roles == nil {
			publicRoutes[paddingRoute] = true
		}
	}
}

/*
In ra danh sách những đường dẫn public không kiểm tra quyền
*/
func DebugPublicRouteRole(app *iris.Application) {
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
	for route, roles := range routesRoles {
		fmt.Println("-" + route)
		for role, allow := range roles {
			if allow.(bool) {
				fmt.Println("      " + RoleName(role)) //allow role in trắng
			} else {
				fmt.Println("     \033[31m^" + RoleName(role) + "\033[0m") //forbid role in đỏ
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
	for path, httpVerbRoles := range pathsRoles {
		fmt.Println("-" + path)
		for httpVerb, roles := range httpVerbRoles {
			fmt.Println("  -" + httpVerb)
			for role := range roles {
				fmt.Println("     " + RoleName(role))
			}
		}
	}
}

/* Insert space between HTTP Verb and Path
Input: GET/blog
Output: GET /blog
*/
func correctRoute(route string) string {
	posFirstSlash := strings.Index(route, "/")
	return route[0:posFirstSlash] + " " + route[posFirstSlash:]
}

//Chuyển role từ in string
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
