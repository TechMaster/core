package router

import (
	"github.com/TechMaster/core/controller"

	"github.com/TechMaster/core/rbac"

	"github.com/kataras/iris/v12"
)

func RegisterRoute(app *iris.Application) {
	app.Get("/", controller.ShowHomePage)
	app.Post("/login", controller.Login)
	rbac.Get(app, "/secret", rbac.AllowAll(), controller.ShowSecret)
	rbac.Get(app, "/logout", rbac.AllowAll(), controller.LogoutFromWeb)
	rbac.Get(app, "/changerole", rbac.Allow(rbac.ADMIN), controller.ShowChangeRoleForm)
	rbac.Post(app, "/changerole", rbac.Allow(rbac.ADMIN), controller.ChangeRole)

	api := app.Party("/api")
	{
		api.Post("/login", controller.LoginREST)
		api.Get("/logout", controller.LogoutREST)
	}
}
