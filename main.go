package main

import (
	"github.com/TechMaster/core/config"
	"github.com/TechMaster/core/controller"
	"github.com/TechMaster/core/email"
	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/router"
	"github.com/TechMaster/core/session"
	"github.com/TechMaster/core/template"
	"github.com/TechMaster/eris"

	"github.com/iris-contrib/middleware/cors"

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

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:9001"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	app.Use(session.Sess.Handler())

	// Load các roles vào bộ nhớ
	rbac.LoadRoles(func() []pmodel.Role {
		return controller.Roles
	})

	rbacConfig := rbac.NewConfig()
	rbacConfig.MakeUnassignedRoutePublic = true
	rbac.Init(rbacConfig) //Khởi động với cấu hình mặc định
	//đặt hàm này trên các hàm đăng ký route - controller
	app.Use(rbac.CheckRoutePermission)

	router.RegisterRoute(app)

	template.InitBlockEngine(app, "./views", "default")
	logger.Log2(eris.SysError("Cuộc sống tốt đẹp. Hãy trận trọng và sống hết mình từng giây một"))
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
	//Khởi động email redis

	asynClient := email.InitRedisMail()
	email.SetDefaultEmailLayout("email_layout") //Set layout mặc định cho các HTML email
	defer asynClient.Close()

	_ = app.Listen(viper.GetString("port"))
}
