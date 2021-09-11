package router

import (
	"github.com/TechMaster/core/controller"
	"github.com/TechMaster/core/ratelimit"
	"github.com/TechMaster/core/rbac"
	"github.com/didip/tollbooth/v6"

	"github.com/kataras/iris/v12"
)

func RegisterRoute(app *iris.Application) {
	limiter := tollbooth.NewLimiter(1, nil) //Maxium 1 request per second

	app.Get("/", controller.ShowHomePage)

	app.Get("/err", controller.ShowErr)

	app.Post("/login", ratelimit.LimitHandler(limiter), controller.Login)
	rbac.Get(app, "/secret", rbac.AllowAll(), controller.ShowSecret)
	rbac.Get(app, "/logout", rbac.AllowAll(), controller.LogoutFromWeb)

	rbac.Get(app, "/changerole", rbac.Allow(rbac.ADMIN), controller.ShowChangeRoleForm)
	rbac.Post(app, "/changerole", rbac.Allow(rbac.ADMIN), controller.ChangeRole)

	api := app.Party("/api")
	{
		api.Post("/login", controller.LoginREST)
		api.Get("/logout", controller.LogoutREST)
		rbac.Get(api, "/books", rbac.Allow(rbac.STUDENT, rbac.TRAINER), controller.Books)
	}
}
