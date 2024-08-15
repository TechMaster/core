package controller

import (
	"strconv"
	"strings"

	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/kataras/iris/v12"
)


func GetAllRule(ctx iris.Context) {
	ctx.ViewData("rules", RulesDb)
	_ = ctx.View("rules")
}

func AddRule(ctx iris.Context) {
	var request struct {
		ID         int
		Name       string
		Roles      string
		AccessType string
		Method     string
		Path       string
		IsPrivate  bool
	}
	if err := ctx.ReadForm(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err.Error()})
		return
	}
	request.ID = RulesDb[len(RulesDb)-1].ID + 1
	arrRoleStr := strings.Split(request.Roles, ",")
	var roles = make([]int, 0, len(arrRoleStr))
	for _, roleStr := range arrRoleStr {
		roleInt, err := strconv.Atoi(roleStr)
		if err == nil {
			roles = append(roles, roleInt)
		}
	}
	rule := pmodel.Rule{
		ID:         request.ID,
		Name:       request.Name,
		Roles:      roles,
		AccessType: request.AccessType,
		Method:     request.Method,
		Path:       request.Path,
		IsPrivate:  request.IsPrivate,
	}
	RulesDb = append(RulesDb, &rule)
	rbac.LoadRules(func() []pmodel.Rule {
		rules := make([]pmodel.Rule, 0)
		for _, rule := range RulesDb {
			rules = append(rules, *rule)
		}
		return rules
	})
	ctx.Redirect("/rules")
}

func ShowViewRuleEdit(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	for _, rule := range RulesDb {
		if rule.ID == id {
			roles := convertSliceIntToString(rule.Roles)
			rolesStr := strings.Join(roles, ",")
			ctx.ViewData("rule", iris.Map{
				"ID":         rule.ID,
				"Name":       rule.Name,
				"Roles":      rolesStr,
				"AccessType": rule.AccessType,
				"Method":     rule.Method,
				"Path":       rule.Path,
				"IsPrivate":  rule.IsPrivate,
			})
			break
		}
	}
	ctx.View("rule-edit")
}

func convertSliceIntToString(intSlice []int) []string {
	strSlice := make([]string, len(intSlice))
	for i, v := range intSlice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strSlice
}

func EditRule(ctx iris.Context) {
	var request struct {
		ID         int
		Name       string
		Roles      string
		AccessType string
		Method     string
		Path       string
		IsPrivate  bool
	}
	if err := ctx.ReadForm(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err.Error()})
		return
	}
	id, _ := ctx.Params().GetInt("id")
	arrRoleStr := strings.Split(request.Roles, ",")
	var roles = make([]int, 0, len(arrRoleStr))
	for _, roleStr := range arrRoleStr {
		roleInt, err := strconv.Atoi(roleStr)
		if err == nil {
			roles = append(roles, roleInt)
		}
	}
	for _, rule := range RulesDb {
		if rule.ID == id {
			rule.Name = request.Name
			rule.Roles = roles
			rule.AccessType = request.AccessType
			rule.Method = request.Method
			rule.Path = request.Path
			rule.IsPrivate = request.IsPrivate
			break
		}
	}
	rbac.LoadRules(func() []pmodel.Rule {
		rules := make([]pmodel.Rule, 0)
		for _, rule := range RulesDb {
			rules = append(rules, *rule)
		}
		return rules
	})
	ctx.Redirect("/rules")
}