# Những thay đổi

### 0.1.21
Sửa lại template để hỗ trợ khởi tạo HTML và Block Engine
```go
func InitHTMLEngine(app *iris.Application, viewFolder string, defaultLayout string) 
func InitBlockEngine(app *iris.Application, viewFolder string, defaultLayout string)
```
Tôi đã copy code ở [https://github.com/kataras/blocks](https://github.com/kataras/blocks) vào thư mục [blocks](blocks)
[Sửa lại lỗi không đặt được layout mặc định](https://github.com/kataras/blocks/issues/2)
### 0.1.19
Sửa lại hàm để hỗ trợ 4 hình thức
1. Debug app
2. Run app bằng `go run main.go`
3. Debug test
4. Run test
   
```go
func IsAppInDebugMode() bool {
	appCommand := os.Args[0]
	if strings.Contains(appCommand, "debug") || //debug ứng dụng trong vscode
		strings.Contains(appCommand, "exe") || //go run main.go
		strings.Contains(appCommand, "go-build") { //run test
		return true
	}
	return false
}
```
### 0.1.16
Trong package session, bỏ `func IsLogin(ctx iris.Context)`, từ nay hãy dùng 2 hàm này
để lấy thông tin người dùng đăng nhập. Ưu tiên hàm `GetAuthInfo` hơn nhé.
```go
func GetAuthInfo(ctx iris.Context) (authinfo *pmodel.AuthenInfo) 
func GetAuthInfoSession(ctx iris.Context) (authinfo *pmodel.AuthenInfo)
```

### 0.1.14
- Sửa lỗi ở hàm `func assignRoles(method string, path string, roles pmodel.Roles)` bằng cách thay các // bằng /
- Cập nhật lại ReadMe.md
- Bỏ role Root, vì lý do bảo mật

### 0.1.13
- Bổ xung chức năng đồng bộ Role trên nhiều thiết bị
- Copy package session từ [https://github.com/kataras/iris/tree/master/sessions](https://github.com/kataras/iris/tree/master/sessions) vào thư mục [sessions](../sessions)
- Chuyển [https://github.com/techmaster/logger](https://github.com/techmaster/logger) vào thư mục [logger](../logger)
- Bổ xung package [pass](../pass) chuyên để băm và so sánh password


### 0.1.12
- Chi tiết hoá struct User trong [/pmodel/user.go](../pmodel/user.go)