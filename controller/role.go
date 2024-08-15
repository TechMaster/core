package controller

import (
	"regexp"

	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/repo"
	"github.com/kataras/iris/v12"
)

var Roles []pmodel.Role = []pmodel.Role{ // fake data được kéo từ database
	{ID: 1, Name: "admin"},
	{ID: 2, Name: "student"},
	{ID: 3, Name: "trainer"},
	{ID: 4, Name: "sale"},
	{ID: 5, Name: "employer"},
	{ID: 6, Name: "author"},
	{ID: 7, Name: "editor"},
	{ID: 8, Name: "maintainer"},
}

var RulesDb []*pmodel.Rule = []*pmodel.Rule{ // fake data được kéo từ database
	{
		ID:         1,
		Name:       "Home",
		Method:     "GET",
		Path:       "/",
		AccessType: "allow_all",
		IsPrivate:  false,
	},
	{
		ID:         2,
		Name:       "Api Books",
		Method:     "GET",
		Path:       "/api/books",
		Roles:      []int{2, 3},
		AccessType: "allow",
		IsPrivate:  true,
	},
}

func GetAllRole(ctx iris.Context) {
	ctx.ViewData("roles", Roles)
	_ = ctx.View("roles")
}

func AddRole(ctx iris.Context) {
	role := pmodel.Role{}
	if err := ctx.ReadForm(&role); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err.Error()})
		return
	}
	role.Name = regexp.MustCompile(`\s+`).ReplaceAllString(role.Name, "_")
	role.ID = Roles[len(Roles)-1].ID + 1
	Roles = append(Roles, role)
	ctx.Redirect("/roles")
}

func DeleteRole(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	for i, role := range Roles {
		if role.ID == id {
			Roles = append(Roles[:i], Roles[i+1:]...)
			break
		}
	}
	for _, user := range repo.Users {
		for i, role := range user.Roles {
			if role == id {
				user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
				break
			}
		}
	}
	for _, rule := range RulesDb {
		for i, role := range rule.Roles {
			if role == id {
				rule.Roles = append(rule.Roles[:i], rule.Roles[i+1:]...)
			}
		}
	}
	rbac.LoadRoles(func() []pmodel.Role {
		return Roles
	})
	rbac.LoadRules(func() []pmodel.Rule {
		rules := make([]pmodel.Rule, 0)
		for _, rule := range RulesDb {
			rules = append(rules, *rule)
		}
		return rules
	})
	ctx.Redirect("/roles")
}
