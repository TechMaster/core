# Một số package hay dùng

[**Hãy đọc kỹ những thay đổi theo phiên bản**](doc/change.md)
### Lộ trình phát triển mới
1. Bổ xung kết nối WebSocket
2. Các hàm bổ trợ cho go-pg
3. Thêm hàm Delete, Put cho resto
4. Thêm hàm dùng chung cho template
5. Thêm chức năng bổ xung REST API để các dịch vụ khác truy vấn lấy thông tin: tên dịch vụ, thời gian kể từ lúc khởi động, danh sách các route + phần quyền...
6. Session cần kéo dài expire date khi người dùng tiếp tục truy cập
7. Bổ xung cách khác để gửi email, gom vào một dịch vụ chuyên biệt để gửi email
8. Chức năng chạy schedule task để dọn dẹp ví dụ dọn thư mục log, xoá bớt orphan entry trong Redis

### Giới thiệu
Module này tổng hợp nhiều package hữu dụng, sử dụng cùng với Iris framework để tạo ra một ứng dụng hoàn chỉnh
1. config: cấu hình, sử dụng [Viper config](https://github.com/spf13/viper)
2. template: chuyên xử lý template engine
3. session: quản lý session, kết nối vào redis. Phụ thuộc vào [iris session](https://github.com/kataras/iris/blob/master/sessions/sessions.go)
4. resto: thư viện rest client hỗ trợ retry, sử dụng [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp)
5. rbac: phân quyền Role Based Access Control
6. pmodel: định nghĩa cấu trúc dữ liệu dùng chung giữa package rbac, session
7. db: kết nối CSDL Postgresql
8. email: gửi email theo nhiều cách khác nhau
9. logger

![](doc/diagram.jpg)

## 1. Hướng dẫn cài đặt module core
```
go get -u github.com/TechMaster/core@main
```

Chú ý do module core luôn đi cùng với module iris, viper do đó bạn cần bổ xung
```
go get -u github.com/kataras/iris/v12@master
```

#### Chạy thử được luôn ví dụ về module core
Để thử nghiệm nhanh các tính năng của module core bằng cách:
```
git clone https://github.com/TechMaster/core.git
cd core
go mod tidy
```

Khởi động một redis server. Trước đó hãy tạo thư mục data để map volume
```
docker run --name=redis -p 6379:6379 -d -e REDIS_PASSWORD=123 -v $PWD/data:/data redis:alpine3.14 /bin/sh -c 'redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}'
```

Chạy lệnh
```
go run main.go
```

Truy cập địa chỉ http://localhost:9001, login thử với các user khác nhau
## 2. Ví dụ hàm main.go sử dụng module core
```go
package main

import (
	"video/router"

	"github.com/TechMaster/core/config"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/session"
	"github.com/TechMaster/core/sessions"
	"github.com/TechMaster/core/tmpl"
	"github.com/TechMaster/core/logger"
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

	// Load các roles vào bộ nhớ
	rbac.LoadRoles(func() []pmodel.Role {
		return controller.Roles
	})

	rbacConfig := rbac.NewConfig()
	rbacConfig.MakeUnassignedRoutePublic = true //mọi route không dùng rbac coi là public
	rbac.Init(rbacConfig) //Khởi động với cấu hình mặc định

	//đặt hàm này trên các hàm đăng ký route - controller
	app.Use(rbac.CheckRoutePermission)

	app.HandleDir("/", iris.Dir("./static"))  //phục vụ thư mục file tĩnh

	router.RegisterRoute(app)  //Cấu hình đường dẫn đến các controller

	template.InitViewEngine(app) //Khởi tạo View Template Engine

	// Meger các route load từ database vào RBAC
	rbac.LoadRules(func() []pmodel.Rule {
		rules := make([]pmodel.Rule, 0)
		for _, rule := range controller.RulesDb {
			rules = append(rules, *rule)
		}
		return rules
	})

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

### 3.1 Hỗ trợ đọc Docker Secret
Docker Secret là một tính năng trong Docker Swarm dùng để lưu thông tin nhạy cảm như password. Như vậy bạn sẽ không để lộ password trong file cấu hình hay file docker-compose.yml.

```yaml
version: "3.8"

secrets:
  pgpass: # Khai báo import pgpass từ bên ngoài
    external: true

services:
  whoami:
    image: main
    secrets:
      - pgpass  # Sử dụng trong dịch vụ này
```

Tiếp đến cấu hình file configure, đánh dấu key `pgpass` bởi string `@@pgpass`. **Nhớ phải có `@@` trước tên secret key **để package config khi đọc sẽ tìm đến file ở thư mục /run/secrets đọc nội dung thực của secret
```json
{
	"database": {
			"user": "root",
			"pgpass": "@@pgpass", 
```


## 4. Logger
**Chú ý: nếu đã dùng module github.com/TechMaster/core thì buộc phải dùng cả github.com/TechMaster/core/logger. Tuyệt đối không dùng github.com/TechMaster/logger !**

Khởi tạo mặc định, bạn cần tạo một thư mục logs để logger ghi vào file log. 
Trong thư mục views cần phải có file [views/error.html](views/error.html) để hiển thị lỗi và [views/info.html](views/info.html) để hiển thị thông tin ra trình duyệt kiểu server side rendering

```go
logFile := logger.Init() //Cần phải có 2 file error.html và info.html ở /views
if logFile != nil {
	defer logFile.Close()
}
```

Khởi tạo có cấu hình riêng. Nếu stack trace lỗi quá dài dòng, bạn có thể đặt `LogConfig.Top` giảm xuống 2 hoặc 3
```go
logFile := logger.Init(LogConfig{
			LogFolder:     "mylog/", // thư mục chứa log file. Nếu rỗng có nghĩa là không ghi log ra file
			ErrorTemplate: "error", // tên view template sẽ render error page
			InfoTemplate:  "info",  // tên view template sẽ render info page
			Top:           4,       // số dòng đầu tiên trong stack trace sẽ được giữ lại
		}
)
if logFile != nil {
	defer logFile.Close()
}
```

#### logger cung cấp 2 hàm log lỗi và 1 hàm log info
```go
func Log(ctx iris.Context, err error)  //Sử dụng trong logic xử lý HTTP request đến Iris framework
func Log2(err error)  //Log lỗi mà không báo về cho HTTP client
func Info(ctx iris.Context, msg string, redirectLink ...string) //Thông báo
```

#### Quy ước báo lỗi như sau:
1. Những lỗi do người dùng (HTTP client) tạo ra không ảnh hưởng gì đến ứng dụng ví dụ nhập sai dữ liệu... chỉ tạo 
```go
logger.Log(ctx, eris.Warning("Email không hợp lệ").BadRequest()
logger.Log(ctx, eris.NewFromMsg(err, "Không đọc được dữ liệu gửi lên").SetType(eris.WARNING).BadRequest())
```
Việc đặt cấp độ lỗi cao chỉ làm rác console log hoặc log file của ứng dụng, không giải quyết được vấn đề gì cả

2. Những lỗi thực sự do ứng dụng gây ra thì đặt ở cấp độ `ERROR`, `SYSERROR`, `PANIC`. Trước khi gọi lệnh panic hãy cố gắng log lỗi ra file

```go
logger.Log2(eris.NewFromMsg(err, "Lỗi đặc biệt nghiêm trọng").SetType(eris.PANIC))  //log ra file
panic(err) //rồi cho hệ thống sập !
```
3. Cần đặt đúng loại lỗi. Việc này giúp client xử lý lỗi tốt hơn như:
```go
eris.Warning().BadRequest() 					//Yêu cầu gửi lên bị lỗi
eris.Warning().UnAuthorized() 				//Người dùng không đủ quyền để thực thi
eris.Warning().NotFound() 						//Không tìm thấy bản ghi hay tài nguyên trên server
eris.Warning().InternalServerError() 	//Dành cho hầu hết lỗi phát sinh phía server
```

## 5. Sử dụng Session
### 5.1 Chạy ứng dụng đơn lẻ độc lập
Nếu bạn viết ứng dụng đơn lẻ thì có thể lưu trực tiếp session vào vùng nhớ của ứng dụng web. Khi này bạn không cần dùng Redis hay bất kỳ CSDL.

Hàm khởi tạo Session trong file main.go sẽ như sau
```go
app.Use(session.Sess.Handler())
```
### 5.2 Nhiều ứng dụng dùng chung session database
Khi có nhiều ứng dụng web, microservice dùng chung một domain nhưng định địa chỉ bằng các sub domain khác nhau, để có được chức năng Single Sign On (đăng nhập một lần, nhưng truy cập được nhiều site cùng chung domain), chúng ta buộc phải lưu session ra database chung ví dụ như Redis.

```go
redisDb := session.InitRedisSession()
defer redisDb.Close()
app.Use(session.Sess.Handler())
```
### 5.3 Làm thế nào để biết người dùng đã đăng nhập?
package session cung cấp 2 hàm

```go
func GetAuthInfo(ctx iris.Context) (authinfo *pmodel.AuthenInfo)
func GetAuthInfoSession(ctx iris.Context) (authinfo *pmodel.AuthenInfo)
```
`GetAuthInfo` đầu tiên sẽ lấy thông tin đăng nhập của người dùng từ `ViewData["authinfo"]` nếu không thấy sẽ tiếp tục gọi vào `GetAuthInfoSession` để lấy thông tin từ session.

Nếu trả về `nil` có nghĩa người dùng chưa đăng nhập. Nếu khác `nil` thì cấu trúc dữ liệu trả về như sau:
```go
type AuthenInfo struct {
	Id       string //unique id của user
	FullName string //họ và tên đầy đủ của user
	Email    string //email cũng phải unique
	Avatar   string //unique id hoặc tên file ảnh đại diện
	Roles    Roles  //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
```

Tôi đã bỏ hàm `func IsLogin(ctx iris.Context)` vì hàm này không trả về đầy đủ được thông tin. Ngược lại hàm `func GetAuthInfo` trả về được id, full name, email, avatar và danh sách roles của người dùng.

Bạn cần lấy danh sách roles của người dùng mảng các chuỗi mô tả role. Hãy truyền `authinfo.Roles` vào hàm này
```go
rbac.RolesNames(roles pmodel.Roles)[]string
```
Bạn cần in ra danh sách role vừa có giá trị int và có chuỗi mô tả để debug cho thuận tiện `3:trainer, 8:maintainer`. Hãy tham khảo hàm [func GetAll](repo/repo.go)
```go
rolesString := ""
for i, role := range user.Roles {
	rolesString += fmt.Sprintf("%d:%s", role, rbac.RoleName(role))
	if i < len(user.Roles)-1 {
		rolesString += ", "
	}
}
```

Bạn có một mảng ```[]int``` thể hiện các role, cần chuyển sang kiểu `type Roles map[int]interface{}`. Hãy dùng [IntArrToRoles](pmodel/role.go)
```go
func IntArrToRoles(intArr []int) Roles 
```
### 5.4 Khi logout bắt buộc phải dùng hàm session.Logout
Hàm này thực hiện việc xoá session id và phần tử session trong tập user.Id. Nó có tác dụng loại bỏ bớt rác trong redis database.
```go
func Logout(ctx iris.Context) error
```

### 5.5 Chức năng cập nhật role chỉ dành cho Admin
Khi Admin thay đổi role người dùng. Người này không cần logout mà role có tác dụng ngay, trên mọi thiết bị anh ta đang đăng nhập.

```go
func UpdateRole(userID string, roles []int) error
```

Xem chi tiết [controller/changerole.go](controller/changerole.go) và [session/update_role.go](session/update_role.go)

Để thực hiện được tính năng này phải lưu quan hệ một user.Id chứa một tập các Session.id. Khi cập nhật Roles cho một user.Id chúng ta nhanh chóng tìm được tất cả các Session của user đó để cập nhật. Ngoài ra phải đặt Expire time để xoá bản ghi này.

Nếu vì một nguyên nhân nào đó, thuật toán đồng bộ Role của user trên mọi thiết bị bị lỗi. User có thể logout rồi login lại.
Thuật toán này chưa hoàn hảo, nó có thể để lại rác trong Redis trong một số trường hợp.

Hiện nay tôi copy toàn bộ package https://github.com/kataras/iris/tree/master/sessions vào thư mục sessions. Hiện chưa sửa đổi gì. Tuy nhiên sẽ fix bug ngay nếu package này có lỗi.
### 5.6 Chức năng Logout
Trong framework Iris, khi người dùng logout ở một trình duyệt trên một thiết bị, không làm sao xoá được key = sessionID. Hàm này không những xoá key = sessionID mà còn sửa lại entry UserID bỏ bớt phần tử sessionID
```go
func Logout(ctx iris.Context) error
```

## 6. Sử dụng RBAC
Cần khởi tạo và cấu hình RBAC trong file main.go
Sau đó trong router viết hàm đăng ký route + roleExp + isPrivate + controller
RBAC hỗ trợ 5 hàm:
1. `Allow(rbac.RoleX, rbac.RoleY)`: cho phép RoleX và RoleY
2. `AllowAll()`: cho phép tất cả các role
4. `Forbid(rbac.RoleA, rbac.RoleB)`: cấm role RoleA, RoleB, các role khác đều được phép
5. `ForbidAll()`: cấm tất cả các role trừ admin

Có thể chuyển app và hoặc đối tượng party vào tham số đầu tiên của rbac

```go
func RegisterRoute(app *iris.Application) {
	// Tất cả route phải được viết qua rbac để kiểm soát
	// Nếu không viết vào rbac nó sẽ là public cho tất cả được truy cập và sẽ không dynamic
	rbac.Get(app, "/", rbac.AllowAll(), false, controller.ShowHomePage)
	rbac.Post("/login", rbac.AllowAll(), false, controller.Login)
	rbac.Get(app, "/logout", rbac.AllowAll(), true,controller.LogoutFromWeb)

	blog := app.Party("/blog")
	{
		rbac.Get(blog, "/", rbac.AllowAll(), false, controller.GetAllPosts)
		rbac.Get(blog, "/all", rbac.AllowAll(), true controller.GetAllPosts)
		rbac.Get(blog, "/create", rbac.ForbidAll(), true, controller.GetAllPosts)
		rbac.Get(blog, "/{id:int}", rbac.ForbidAll(), true,controller.GetPostByID)
		rbac.Get(blog, "/delete/{id:int}", rbac.ForbidAll(), true, controller.DeletePostByID)
		rbac.Any(blog, "/any", rbac.ForbidAll(), true, controller.PostMiddleware)
	}

	student := app.Party("/student")
	{
		rbac.Get(student, "/submithomework", rbac.ForbidAll(), true, controller.SubmitHomework)
	}

	trainer := app.Party("/trainer")
	{
		rbac.Get(trainer, "/createlesson", rbac.ForbidAll(), true, controller.CreateLesson)
	}

	sysop := app.Party("/sysop")
	{
		rbac.Get(sysop, "/backupdb", rbac.ForbidAll(), true, controller.BackupDB)
		rbac.Get(sysop, "/upload", rbac.ForbidAll(), true, controller.ShowUploadForm)
		rbac.Post(sysop, "/upload", rbac.ForbidAll(), true, iris.LimitRequestBodySize(300000), controller.UploadPhoto)
	}
}
```

Sau đó sẽ gọi `rbac.LoadRules` để load tất cả các rules đã được phân quyền từ database
```go

rbac.LoadRules(func() []pmodel.Rule {
	// Viết SQL lấy các rules từ database về
	return rules
})
```

Cuối cùng gọi hàm `BuildPublicRoute` để phân các route public

```go
//Luôn để hàm này sau tất cả lệnh cấu hình đường dẫn với RBAC
rbac.BuildPublicRoute(app)
```

## 7. Cấu trúc dữ liệu trong pmodel

pmodel là nơi định nghĩa cấu trúc dữ liệu phụ vụ việc đăng nhập, quản lý người dùng

Cấu trúc chi tiết User xem ở đây [pmodel/user.go](pmodel/user.go)

Danh sách các Role cấp cho một user
```go
type Roles map[int]interface{}
```
Struct sẽ lưu trong session để hệ thống quản lý phiên đăng nhập của người dùng
```go
type AuthenInfo struct {
	Id       string //unique id của user
	FullName string //họ và tên đầy đủ của user
	Email    string //email cũng phải unique
	Avatar   string //unique id hoặc tên file ảnh đại diện
	Roles    Roles  //kiểu map[int]bool. Cần phải chuyển đổi Roles []int32 `pg:",array"` sang
}
```

Chú ý kiểu `map[int]bool` khi lưu vào Redis sẽ biến thành `map[string]bool`
Chuyển đổi Roles kiểu map[int] bool sang mảng []int để lưu xuống CSDL
```go
func RolesToIntArr(roles Roles) []int
```

Chuyển đổi kiểu intArray trong đó mỗi phần tử ứng với một role, sang kiểu map[int] bool
```go
func IntArrToRoles(intArr []int) Roles
```

Cấu trúc chi tiết Rule xem ở đây [pmodel/rule.go](pmodel/rule.go)

Danh sách các rules sẽ được lấy từ database

```go
/*
Rule là cấu hình cho việc kiểm tra quyền truy cập
- Nếu IsPrivate = true thì cần kiểm tra quyền và Roles, AccessType không có ý nghĩa, có thể để AccessType = "allow_all" cho dễ hiểu
*/
type Rule struct {
	ID         int    //ID của rule
	Name       string //Tên của rule
	Roles      []int  `pg:",array"` //Danh sách các role có thể truy xuất
	AccessType string //Allow, AllowAll, ForbidAll, Forbid, ForbidAll
	Method     string //GET, POST, PUT, DELETE, PATCH
	Path       string //Đường dẫn
	IsPrivate  bool   //true: cần kiểm tra quyền, false: không cần kiểm tra quyền
	Services   string //Dịnh nghĩa các rule cho các service khác nhau
}
```

## 8. Template Engine
Trong template có 2 hàm khởi tạo View Template Engine. 
- `viewFolder`: thư mục chứa view template. Mặc định nên để là [views](views)
- `defaultLayout`: tên file template layout mặc định. Mặc định nó phải nằm trong thư mục `layouts` bên trong thư mục chứa view template.

```go
//Nếu bạn dùng HTML template
func InitHTMLEngine(app *iris.Application, viewFolder string, defaultLayout string)

//Nếu bạn dùng Block template
func InitBlockEngine(app *iris.Application, viewFolder string, defaultLayout string)
```

Tốt nhất bạn nên tuân thủ cấu trúc thư mục như sau
```
views
├── layouts
│   └── default.html  -> layout mặc định
├── partials  -> template component hiển thị một phần hay dùng lại của trang web
│   ├── footer.html
│   ├── header.html
│   └── menu.html
├── changerole.html
├── error.html -> trang hiển thị lỗi. logger.Log cần phải có
├── index.html
└── info.html -> trang hiển thị thông báo. logger.Log cần phải có
```


Tôi đã copy code ở [https://github.com/kataras/blocks](https://github.com/kataras/blocks) vào thư mục [blocks](blocks)
[Sửa lại lỗi không đặt được layout mặc định](https://github.com/kataras/blocks/issues/2)

Trong hàm [main.go](main.go) cấu hình view template engine như sau

```go
template.InitBlockEngine(app, "./views", "default")
```

Hỏi: Làm sao để bổ xung custom function vào view template engine?

Đáp: Hãy gọi biến global ứng với template engine bạn chọn. Ví dụ bạn dùng BlockEngine
```go
template.BlockEngine.AddFunc("listmenu", func() []string {
	return []string{"home", "products", "about"}
})
```



## 9. Resto thư viện REST client dựa trên cơ chế retry
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

## 10. db kết nối CSDL Postgresql
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

## 11. email

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

Phiên bản 0.1.28 bổ xung thêm gửi email qua Asynq và Redis Stream
Xem [redis_mail.go](email/redis_mail.go)

## 12. pass các hàm băm password
Tuyệt đối không được lưu secret key hay các chuỗi nhạy cảm vào đây. Xem chi tiết [pass/password.go](pass/password.go)
Tạo ra một interface chung cho tất cả các thư viện băm password tuân thủ
```go
type PasswordLib interface {
	Hash(password string) (hash string)
	Compare(password string, hashpass string) bool
}
```

Hiện tại dùng thư viện [Argon2id](https://pkg.go.dev/golang.org/x/crypto/argon2) là thư viện password tốt nhất hiện nay

### 12.1 Băm password

```go
func HashPassword(inputpass string) string {
	return PassLib.Hash(inputpass)
}
```

### 12.2 Kiểm tra password
```go
/*
Hàm check password hỗ trợ cả kiểu SHA1 cũ và bcrypt mới
- inputpass: password nhập vào lúc login
- hashedpass: password đã băm đọc từ CSDL
- salt: chuỗi nhiễu tạo ra từ thuật toán SHA1 cũ
*/
func CheckPassword(inputpass string, hashedpass string, salt string) bool {
	if salt != "" {
		pass := p.NewPassword(sha1.New, 50, 64, 10000)
		return pass.VerifyPassword(inputpass, hashedpass, salt) //Sửa theo yêu cầu Nhật Đức
	} else {
		return PassLib.Compare(inputpass, hashedpass)
	}
}
```

## 13. Rate Limit
Sử dụng `github.com/didip/tollbooth/v6`. Tolbooth là middleware để giới hạn số lượng http request đến trong một khoảng thời gian. Cấu hình trong router như sau:

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

## Để phát hành phiên bản mới module này cần làm những bước sau
Thay v0.1.3 bằng phiên bản thực tế
```
git add .
git commit -m "v0.1.10"
git tag v0.1.10
git push origin v0.1.10
GOPROXY=proxy.golang.org go list -m github.com/TechMaster/core@v0.1.10
```
