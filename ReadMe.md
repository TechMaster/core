# Một số package hay dùng

Module này tổng hợp nhiều package hữu dụng, sử dụng cùng với Iris framework để tạo ra một ứng dụng hoàn chỉnh
1. config: cấu hình, sử dụng [Viper config](https://github.com/spf13/viper)
2. template: chuyên xử lý template engine
3. session: quản lý session, kết nối vào redis. Phụ thuộc vào [iris session](https://github.com/kataras/iris/blob/master/sessions/sessions.go)
4. resto: thư viện rest client hỗ trợ retry, sử dụng [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp)
5. rbac: phân quyền Role Based Access Control
6. pmodel: định nghĩa cấu trúc dữ liệu dùng chung giữa package rbac, session
7. db: kết nối CSDL Postgresql
8. email: gửi email theo nhiều cách khác nhau

![](doc/diagram.jpg)

## 1. Hướng dẫn cài đặt module core
```
go get -u github.com/TechMaster/core@main
```

Chú ý do module core luôn đi cùng với module iris, viper do đó bạn cần bổ xung
```
go get -u github.com/kataras/iris/v12@master
```

## 2. Ví dụ hàm main.go sử dụng module core
```go
package main

import (
	"video/router"

	"github.com/TechMaster/core/config"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/session"
	"github.com/TechMaster/core/template"
	"github.com/TechMaster/logger"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

func main() {
	app := iris.New()
	config.ReadConfig()

	logFile := logger.Init() //Cần phải có 2 file error.html và info.html ở /views
	if logFile != nil {
		defer logFile.Close()
	}

	redisDb := session.InitRedisSession()
	defer redisDb.Close()
	app.Use(session.Sess.Handler())

	rbacConfig := rbac.NewConfig()
	rbacConfig.RootAllow = false  //cấm không dùng tài khoản root
	rbacConfig.MakeUnassignedRoutePublic = true //mọi route không dùng rbac coi là public
	rbac.Init(rbacConfig) //Khởi động với cấu hình mặc định

	//đặt hàm này trên các hàm đăng ký route - controller
	app.Use(rbac.CheckRoutePermission)

	app.HandleDir("/", iris.Dir("./static"))  //phục vụ thư mục file tĩnh

	router.RegisterRoute(app)  //Cấu hình đường dẫn đến các controller

	template.InitViewEngine(app) //Khởi tạo View Template Engine

	//Luôn để hàm này sau tất cả lệnh cấu hình đường dẫn với RBAC
	rbac.BuildPublicRoute(app)
	//rbac.DebugRouteRole()
	_ = app.Listen(viper.GetString("port"))
}
```

## 3. Sử dụng config

Cần đảm bảo phải có 2 file `config.dev.json` và `config.product.json` ở thư mục gốc của dự án
```
.
├── config.dev.json
├── config.product.json
```

Để đọc giá trị cấu hình dùng lệnh `viper.GetString("key")` hoặc `viper.GetInt("key")`
```go
_ = app.Listen(viper.GetString("port"))
```

Tham khảo file cấu hình [config.dev.json](config.dev.json)

package config có hàm này để cho biết ứng dụng đang chạy mode Debug hay Production từ đó nạp file cấu hình cho phù hợp
```go
/*
Trả về true nếu ứng dụng đang chạy ở chế độ Debug và ngược lại
*/
func IsAppInDebugMode() bool {
	appCommand := os.Args[0]
	if strings.Contains(appCommand, "debug") || strings.Contains(appCommand, "exe") {
		return true
	}
	return false
}
```

## 4. Sử dụng Session
### 4.1 Chạy ứng dụng đơn lẻ độc lập
Nếu bạn viết ứng dụng đơn lẻ thì có thể lưu trực tiếp session vào vùng nhớ của ứng dụng web. Khi này bạn không cần dùng Redis hay bất kỳ CSDL.

Hàm khởi tạo Session trong file main.go sẽ như sau
```go
session.InitSession()
app.Use(session.Sess.Handler())
```
### 4.2 Nhiều ứng dụng dùng chung session database
Khi có nhiều ứng dụng web, microservice dùng chung một domain nhưng định địa chỉ bằng các sub domain khác nhau, để có được chức năng Single Sign On (đăng nhập một lần, nhưng truy cập được nhiều site cùng chung domain), chúng ta buộc phải lưu session ra database chung ví dụ như Redis.

```go
redisDb := session.InitRedisSession()
defer redisDb.Close()
app.Use(session.Sess.Handler())
```

## 5. Sử dụng RBAC
Cần khởi tạo và cấu hình RBAC trong file main.go
Sau đó trong router viết hàm đăng ký route + controller

RBAC hỗ trợ 4 hàm:
1. `Allow(rbac.RoleX, rbac.RoleY)`: cho phép RoleX và RoleY
2. `AllowAll()`: cho phép tất cả các role
3. `Forbid(rbac.RoleA, rbac.RoleB)`: cấm role RoleA, RoleB, các role khác đều được phép
4. `ForbidAll()`: cấm tất cả các role

Có thể chuyển app và hoặc đối tượng party vào tham số đầu tiên của rbac

```go
func RegisterRoute(app *iris.Application) {

	app.Get("/", controller.ShowHomePage) //Không dùng rbac có nghĩa là public method
	app.Post("/login", controller.Login)
	rbac.Get(app, "/logout", rbac.AllowAll(), controller.LogoutFromWeb)

	blog := app.Party("/blog")
	{
		blog.Get("/", controller.GetAllPosts) //Không dùng rbac có nghĩa là public method
		rbac.Get(blog, "/all", rbac.AllowAll(), controller.GetAllPosts)
		rbac.Get(blog, "/create", rbac.Forbid(rbac.MAINTAINER), controller.GetAllPosts)
		rbac.Get(blog, "/{id:int}", rbac.Allow(rbac.AUTHOR, rbac.EDITOR), controller.GetPostByID)
		rbac.Get(blog, "/delete/{id:int}", rbac.Allow(rbac.ADMIN, rbac.AUTHOR, rbac.EDITOR), controller.DeletePostByID)
		rbac.Any(blog, "/any", rbac.Allow(rbac.MAINTAINER), controller.PostMiddleware)
	}

	student := app.Party("/student")
	{
		rbac.Get(student, "/submithomework", rbac.Allow(rbac.STUDENT), controller.SubmitHomework)
	}

	trainer := app.Party("/trainer")
	{
		rbac.Get(trainer, "/createlesson", rbac.Allow(rbac.TRAINER), controller.CreateLesson)
	}

	sysop := app.Party("/sysop")
	{
		rbac.Get(sysop, "/backupdb", rbac.Allow(rbac.SYSOP), controller.BackupDB)
		rbac.Get(sysop, "/upload", rbac.Allow(rbac.MAINTAINER, rbac.SYSOP), controller.ShowUploadForm)
		rbac.Post(sysop, "/upload", rbac.Allow(rbac.MAINTAINER, rbac.SYSOP, rbac.SALE), iris.LimitRequestBodySize(300000), controller.UploadPhoto)
	}
}
```
Mặc định đã có sẵn các role sau đây

```go
const (
	ROOT       = 0 //Role đặc biệt, vượt qua mọi logic kiểm tra quyền khi config.RootAllow = true.
	ADMIN      = 1
	STUDENT    = 2
	TRAINER    = 3
	SALE       = 4
	EMPLOYER   = 5
	AUTHOR     = 6
	EDITOR     = 7 //edit bài, soạn page, làm công việc digital marketing
	MAINTAINER = 8 //quản trị hệ thống, gánh bớt việc cho Admin, back up dữ liệu. Sửa đổi profile,role user, ngoại trừ role ROOT và Admin
)
```
## 6. Cấu trúc dữ liệu trong pmodel

pmodel là nơi định nghĩa cấu trúc dữ liệu phụ vụ việc đăng nhập, quản lý người dùng

Danh sách các Role cấp cho một user
```go
type Roles map[int]interface{}
```

Thông tin tài khoản tối giản của người dùng
```go
//Thông tin tài khoản
type User struct {
	User  string
	Pass  string
	Email string
	Roles Roles
}
```

Struct sẽ lưu trong session để hệ thống quản lý phiên đăng nhập của người dùng
```go
type AuthenInfo struct {
	User  string
	Email string
	Roles Roles //kiểu map[int]bool
}
```

Chú ý kiểu `map[int]bool` khi lưu vào Redis sẽ biến thành `map[string]bool`

## 7. Template Engine
Hiện chưa viết được nhiều hàm phụ trợ. Sau sẽ bổ xung thêm.
Chủ yếu sử dụng Blocks template của iris. Nếu thư viện này có lỗi sẽ clone và tạo thư viện mới.
Chú ý để dùng được `*view.BlocksEngine` bạn phải lấy bản mới nhất thư viện Iris
```
go get -u github.com/kataras/iris/v12@master
```
[template/base.go](template/base.go)
```go
package template

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/view"
)

var ViewEngine *view.BlocksEngine

func InitViewEngine(app *iris.Application) {
	ViewEngine = iris.Blocks("./views", ".html")
	app.RegisterView(ViewEngine)
}
```
## 8. Resto thư viện REST client dựa trên cơ chế retry
```go
response, err := resto.Retry(numberOfTimesToTry, numberOfMilliSecondsToWait).Post(url, jsondata)
response, err := resto.Retry(numberOfTimesToTry, numberOfMilliSecondsToWait).Get(url)
```

Ví dụ chi tiết
```go
response, err := resto.Retry(5, 1000).Post("http://auth/api/login", loginReq)
if err != nil {
  logger.Log(ctx, eris.NewFromMsg(err, "Lỗi khi gọi Auth service").InternalServerError())
  return
}
if response.StatusCode != iris.StatusOK {
  var res struct {
    Error string `json:"error"`
  }
  _ = json.NewDecoder(response.Body).Decode(&res)
  logger.Log(ctx, eris.Warning(res.Error).UnAuthorized())
  return
}
```

## 9. db kết nối CSDL Postgresql
```go
db.ConnectPostgresqlDB(config.Config) //Kết nối vào  CSDL
defer db.DB.Close()
```

Cấu hình kết nối CSDL để ở trong file `config.dev.json` và `config.product.json`
```json
{
	"database": {
			"user": "postgres",
			"password": "123",
			"database": "iris",
			"address": "localhost:5432"
	},
}
```

## 10. email

Đầu tiên là interface gửi email trong [mail_sender.go](email/mail_sender.go)
```go
type MailSender interface {
	SendPlainEmail(to []string, subject string, body string) error
	SendHTMLEmail(to []string, subject string, tmplFile string, data map[string]interface{}) error
}
```

[gmail_smtp.go](email/gmail_smtp.go) gửi email sử dụng một tài khoản Gmail phải bật chế độ không an toàn mới gửi được. Còn cấu hình bằng OAuth2 Gmail Service thì quá khó. Tôi bó tay.

[fake_gmail.go](email/fake_gmail.go) cũng gửi đi từ một tài khoản Gmail, nhưng địa chỉ thư nhận luôn là một hòm thư cấu hình sẵn `test_receive_email string` dùng để kiểm tra, debug ứng dụng.
```go
func InitFakeGmail(config *SMTPConfig, test_receive_email string)
```

[email_db.go](email/email_db.go) thay vì gửi email thì tạo một records trong bảng `debug.emailstore` CSDL Postgresql. Cấu trúc bảng như dưới.

```go
type EmailStore struct {
	tableName  struct{} `pg:"debug.emailstore"`
	Id         int      `pg:",pk"`
	Receipient string
	Subject    string
	Body       string
	CreatedAt  time.Time
}
```
Trong tương lai tôi sẽ bổ xung thêm vài biến thể gửi mail tuân thủ `type MailSender interface`