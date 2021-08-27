package controller

import (
	"strconv"
	"strings"

	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/core/repo"
	"github.com/TechMaster/core/session"
	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12"
)

func ShowChangeRoleForm(ctx iris.Context) {
	_ = ctx.View("changerole")
}

func ChangeRole(ctx iris.Context) {
	var request struct {
		Email string
		Roles string
	}

	if err := ctx.ReadForm(&request); err != nil {
		logger.Log(ctx, eris.Warning("Cannot read form"))
		return
	}

	user, err := repo.QueryByEmail(strings.ToLower(request.Email))
	if err != nil {
		logger.Log(ctx, eris.Warning("User not found"))
		return
	}

	//Chuyển danh sách role dạng chuỗi "2,4,5,6" thành []int{2,4,5,6}
	arrRoleStr := strings.Split(request.Roles, ",")
	var roles = make([]int, 0, len(arrRoleStr))
	for _, roleStr := range arrRoleStr {
		roleInt, err := strconv.Atoi(roleStr)
		if err == nil {
			roles = append(roles, roleInt)
		}
	}

	//Cập nhật vào database
	user.Roles = roles
	repo.UpSertUser(user) //Cập nhật role trong repo

	//Cập nhật vào Redis Session
	err = session.UpdateRole(user.Id, roles)
	if err != nil {
		logger.Log(ctx, eris.NewFromMsg(err, "Failed to update role"))
		return
	}
	_, _ = ctx.WriteString("Update role successfully")
}
