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

	rbac.Get(app, "/", rbac.AllowAll(), false, controller.ShowHomePage)

	app.Get("/err", controller.ShowErr)

	rbac.Post(app, "/login", rbac.AllowAll(), false, ratelimit.LimitHandler(limiter), controller.Login)
	rbac.Get(app, "/secret", rbac.AllowAll(), true, controller.ShowSecret)
	rbac.Get(app, "/logout", rbac.AllowAll(), true, controller.LogoutFromWeb)

	rbac.Get(app, "/roles", rbac.AllowOnlyAdmin(), true, controller.GetAllRole)
	rbac.Post(app, "/roles", rbac.AllowOnlyAdmin(), true, controller.AddRole)
	rbac.Get(app, "/roles/{id}", rbac.AllowOnlyAdmin(), true, controller.DeleteRole)

	rbac.Get(app, "/rules", rbac.AllowOnlyAdmin(), true, controller.GetAllRule)
	rbac.Post(app, "/rules", rbac.AllowOnlyAdmin(), true, controller.AddRule)
	rbac.Get(app, "/rules/{id}", rbac.AllowOnlyAdmin(), true, controller.ShowViewRuleEdit)
	rbac.Post(app, "/rules/{id}", rbac.AllowOnlyAdmin(), true, controller.EditRule)

	rbac.Get(app, "/changerole", rbac.AllowOnlyAdmin(), true, controller.ShowChangeRoleForm)
	rbac.Post(app, "/changerole", rbac.AllowOnlyAdmin(), true, controller.ChangeRole)

	api := app.Party("/api")
	{
		rbac.Post(api, "/login", rbac.AllowAll(), false, controller.LoginREST)
		rbac.Get(api, "/logout", rbac.AllowAll(), false, controller.LogoutREST)
		rbac.Get(api, "/books", rbac.AllowOnlyAdmin(), true, controller.Books)
	}
}
