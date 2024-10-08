# Những thay đổi

## 0.1.82: 29/8/2024

Chỉnh sửa, nâng cấp RBAC (Role Base Access Control)
- Nâng cấp cho các roles dynamic 
- Cấp quyền truy cập vào route
- Tự động load các roles, rules
- Tự động insert các rules và delete (Nếu không tồn tại các rule đó sẽ xóa trong db)

## 0.1.51: 1/4/2022
- Thay đổi cách log error ra console, thống nhất trả ra kiểu text đối với content-type application/json thay vì lúc json, lúc text. File thay đổi [ở đây](../logger/log_error.go)
- Import module eris mới nhất

## 0.1.50: 22/03/2022
- Bổ sung tính năng nếu chương trình đang chạy ở chế độ debug thì log query ra console. File thay đổi [ở đây](../db/db.go)

## 0.1.46: 09/03/2022
- Bổ sung gửi email marketing. File bổ sung [ở đây](../email/redis_mail.go)


### 0.1.42: 30/11/2021
- Đức thay đổi `UpdateRole` thành `UpdateUserInfo`, từ giờ khi cập nhật propfile hay cập nhật role sẽ đều dùng chung `UpdateUserInfo` này. File thay đổi [ở đây](../session/update_user.go)

### 0.1.40: 28/9/2021
- Bổ sung thêm nếu template engine khác nil thì mới Set viewdata authinfo ở hàm `CheckPermission` rbac.

### 0.1.39: 27/9/2021
Đức [cập nhật GetAuthInfo](../session/query_session.go) đọc SessId từ Header Request để bên frontend có thể lấy được userinfo.
Ngoài ra thay đổi biến redisClient thành biến public RedisClient để ở student có thể lấy được biến này.

### 0.1.37: 17/9/2021
Cường: cập nhật package [config](../config/config.go) để có thể đọc được thông tin nhạy cảm trong Docker Secret khi triển khai lên Docker Swarm.

Hàm `func ParseViperSettings()` trong  [config](../config/config.go) sẽ đọc toàn bộ thông tin cấu hình sau khi viper đọc xong từ file cấu hình sau đó tìm đến những trường dùng Docker Secret để đọc giá trị thực trong file trong thư mục chuyên lưu Docker secret `/run/secrets`

Trong file cấu hình `config.product.json` key nào cần đọc Docker Secret hãy khai báo như sau:

```json
"password": "@@pg_password", 
```

Phần giá trị của string **bổ xung 2 ký tự `@@`** làm chỉ dấu để [config](../config/config.go) khi đọc tới, tiếp tục đọc file lưu trong thư mục `/run/secrets`. Các mẫu giá trị Docker Secret được ghi vào từng file trong thư mục `/run/secrets` khi container khởi động.

Ngoài ra phải cấu hình service trong docker-compose.yml để service dùng cụ thể Docker secret nào
```yaml
version: "3.8"

secrets:
  pg_password: # Khai báo import pg_password từ bên ngoài
    external: true

services:
  whoami:
    image: main
    secrets:
      - pg_password  # Sử dụng trong dịch vụ này
```




### 0.1.34: 17/9/2021
Bổ sung thêm trường hợp nếu lỗi là broken pipe thì ignore nó đi, không ghi ra console để tránh log lỗi không cần thiết quá nhiều. Tham khảo về lỗi broken pipe [tại đây](https://noknow.info/it/go/handling_error_broken_pipe?lang=en).
```go
if errors.Is(err, syscall.EPIPE){
	return
}
```
File bổ sung thêm code [ở đây](../logger/log_error.go)
### 0.1.33: 11/9/2021

**Bổ xung package** [ratelimit.go](../ratelimit/ratelimit.go) giới hạn số request xử lý trong một giây
Xem ví dụ ở [router/base.go](../router/base.go)
```go
import (
	"github.com/TechMaster/core/ratelimit"
	"github.com/didip/tollbooth/v6"
)

func RegisterRoute(app *iris.Application) {
	limiter := tollbooth.NewLimiter(1, nil) //Tối đa 1 request trong 1 giây
	app.Post("/login", ratelimit.LimitHandler(limiter), controller.Login)  //áp dụng với POST /login
}
```

**Sửa đổi package** [pass](../pass/password.go): thay thế hàm băm password từ BCrypt thành [Argon2id](https://pkg.go.dev/golang.org/x/crypto/argon2#hdr-Argon2id)

Hãy dùng 2 hàm chính là
```go
func HashPassword(inputpass string) string
func CheckPassword(inputpass string, hashedpass string, salt string) bool
```

### 0.1.32: 7/9/2021
Thay đổi trong [session/session.go](../session/session.go)
```go
const (
	SESSION_COOKIE = "mycookiesession"
	SESS_USER      = "authenticate"
	AUTHINFO       = "authinfo"
)

//Các biến dùng chung trong packge
var Sess *sessions.Sessions         //Cấu hình Session Manager
var redisDB *redis_session.Database //Đây là một wrapper nối xuống Redis của Iris
var redisClient *redis.Client       //Đây là redis client trực tiếp nối xuống Redis db không qua Iris
var expires = time.Hour * 720         //Thời gian mà 1 session sẽ hết hạn và bị xoá khỏi Redis
```
Thay đổi trong [pmodel/user.go](../pmodel/user.go)
```go
type AuthenInfo struct {
	UserId       string //unique id của user
	UserFullName string //họ và tên đầy đủ của user
	UserEmail    string //email cũng phải unique
	UserAvatar  string //unique id hoặc tên file ảnh đại diện
	UserPhone	string
	Roles        Roles  //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
```
Và một số thay đổi khác nữa trong controller do tên field AuthenInfo thay đổi

### 0.1.30 : 6/9/2021

Thay đổi trong [blocks/engine.go](../blocks/engine.go)
```go
func (s *BlocksEngine) ExecuteWriter(w io.Writer, tmplName, layoutName string, data interface{}) error {
	if layoutName == "" {  //Nếu tham số rỗng, thì dùng defaultLayoutName
		layoutName = s.Engine.defaultLayoutName
	}

	if layoutName == view.NoLayout { //Để muốn không dùng layout thì truyền vào iris.nolayout
		layoutName = ""
	}

	return s.Engine.ExecuteTemplate(w, tmplName, layoutName, data)
}
```
Đã bổ xung hàm kiểm thử ở [blocks/block_test.go](../blocks/block_test.go). Chạy debug test từng hàm.

Có mấy trường hợp:

1. Sử dụng default layout
Default layout sẽ là file [views/layouts/default.html](../views/layouts/default.html)
```go
template.InitBlockEngine(app, "./views", "default")
```

Trong các handler
```go
ctx.View("template")  //Hàm này sẽ default layout
```

2. Sử dụng custom layout
Trong thư mục /views/layouts phải có file custom_layout.html
```go
ctx.ViewLayout("custom_layout")
ctx.View("template")
```

3. Hoàn toàn không dùng layout
```go
ctx.ViewLayout(view.NoLayout)
ctx.View("template")
```


Thay đổi hàm  SendHTMLEmail trong package email
```go
//Cũ 
SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error

//Mới
SendHTMLEmail(to []string, subject string, data map[string]interface{}, tmpl_layout ...string) error
```
Tham số variadic `tmpl_layout` có tối thiểu một tham số template, tham số thứ 2 là layout.
Nếu không truyền vào layout thì sẽ lấy `defaultEmailLayout` được cấu hình qua hàm

```go
func SetDefaultEmailLayout(defaultLayout string) {
	defaultEmailLayout = defaultLayout
}
```


### 0.1.28
Bổ xung tính năng gửi email sử dụng Asynq và Redis Stream
Xem file [redis_mail.go](../email/redis_mail.go)

Kiểm thử ở [redis_mail_test.go](../email/redis_mail_test.go)

Cấu hình ở file [main.go](../main.go)
```go
asynClient := email.InitRedisMail()
defer asynClient.Close()
```
### 0.1.27

Bổ xung ứng dụng Vue3 [bookvue](../bookvue/ReadMe.md) demo tính năng tự động login của Vuejs
[bookvue/src/App.vue](../bookvue/src/App.vue)
```javascript
async fetchBooks() {
	try {
		this.error = null
		this.loading = true
		const url = `http://localhost:9001/api/books`
		axios.defaults.withCredentials = true;
		axios.defaults.headers.post['Content-Type'] = 'application/json';
		const response = await axios.get(url)

		if (response.status != 200) {
			console.log(response)
		} else {
			this.books = response.data.books
			this.authinfo = response.data.authinfo
		}
		
	} catch (err) {       
		console.log(err)
	}
	this.loading = false
}
```

File [main.go](../main.go) bổ xung thêm CORS middleware
```go
crs := cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:9001"},
	AllowCredentials: true,
})
app.UseRouter(crs)
```

File [pass.go](../pass/password.go) sửa lại dòng này thành
```go
return pass.VerifyPassword(inputpass, hashedpass, salt) //Sửa theo yêu cầu Nhật Đức
```

### 0.1.26

Bổ xung hàm `func Log2(err error)` trong package logger. Xem tại [logger/log2.go](../logger/log2.go).
Hàm này cho phép log lỗi mà không cần tham số `ctx iris.Context`.

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