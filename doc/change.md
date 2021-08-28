# Những thay đổi


### 0.1.25
Sửa lỗi khi hàm [CheckRoutePermission](../rbac/check_permission.go) trong package rbac gọi
`logger.Log()` thì logger không tìm được view template để hiển thị lỗi. 

Mặc dù trước đó đã khởi tạo trong hàm `func Init(logConfig ...LogConfig) *os.File`

Thử chuyển `var logConfig LogConfig` sang thành `var LogConf *LogConfig`


### 0.1.24
Fix lỗi Internal Server Error do không có file [views/layout/default.html](../views/layouts/default.html)

Trước khi sửa file [template/base.go](../template/base.go)
```go
func InitViewEngine(app *iris.Application) {
	InitBlockEngine(app, "./views", "default")
}
```
Sau khi sửa, gán tham số default layout là ""
```go
func InitViewEngine(app *iris.Application) {
	InitBlockEngine(app, "./views", "")
}
```
### 0.1.23
Sửa lỗi khi eris.Err cấp độ từ Error, SysError, Panic không in ra console đủ strack track

Nếu khởi tạo giá trị của ErisStringFormat ở scope package global, lúc này `logConfig.Top` chưa được khởi tạo giá trị, mặc định là 0 dẫn đến không in ra được stack track
```go
var ErisStringFormat = eris.StringFormat{
		Options: eris.FormatOptions{
			InvertOutput: false, // flag that inverts the error output (wrap errors shown first)
			WithTrace:    true,  // flag that enables stack trace output
			InvertTrace:  true,  // flag that inverts the stack trace output (top of call stack shown first)
			WithExternal: false,
			Top:          logConfig.Top, // Chỉ lấy 3 dòng lệnh đầu tiên
			//Mục tiêu để báo lỗi gọn hơn, stack trace đủ ngắn
		},
		MsgStackSep:  "\n",  // separator between error messages and stack frame data
		PreStackSep:  "\t",  // separator at the beginning of each stack frame
		StackElemSep: " | ", // separator between elements of each stack frame
		ErrorSep:     "\n",  // separator between each error in the chain
	}
```

Cần phải chuyển vào

```go
var ErisStringFormat eris.StringFormat //khai báo biến global ErisStringFormat

func Init(_logConfig ...LogConfig) *os.File {
	if len(_logConfig) > 0 {
		logConfig = _logConfig[0]
	} else { //Truyền cấu hình nil thì tạo cấu hình mặc định
		logConfig = LogConfig{
			LogFolder:     "logs/", // thư mục chứa log file. Nếu rỗng có nghĩa là không ghi log ra file
			ErrorTemplate: "error", // tên view template sẽ render error page
			InfoTemplate:  "info",  // tên view template sẽ render info page
			Top:           3,       // số dòng đầu tiên trong stack trace sẽ được giữ lại
		}
	}

	//Khởi tạo biến ở đây, sau khi logConfig.Top được gán giá trị mới đúng
	ErisStringFormat = eris.StringFormat{
		Options: eris.FormatOptions{
			InvertOutput: false, // flag that inverts the error output (wrap errors shown first)
			WithTrace:    true,  // flag that enables stack trace output
			InvertTrace:  true,  // flag that inverts the stack trace output (top of call stack shown first)
			WithExternal: false,
			Top:          logConfig.Top, // Chỉ lấy 3 dòng lệnh đầu tiên
			//Mục tiêu để báo lỗi gọn hơn, stack trace đủ ngắn
		},
		MsgStackSep:  "\n",  // separator between error messages and stack frame data
		PreStackSep:  "\t",  // separator at the beginning of each stack frame
		StackElemSep: " | ", // separator between elements of each stack frame
		ErrorSep:     "\n",  // separator between each error in the chain
	}
}
```

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