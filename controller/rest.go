package controller

import (
	"github.com/TechMaster/core/repo"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/session"

	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12"
)

/*
Login thông qua axios.post dành cho ứng dụng Vue
Request.ContentType = 'application/json'
*/
func LoginREST(ctx iris.Context) {
	var loginReq LoginRequest

	if err := ctx.ReadJSON(&loginReq); err != nil {
		logger.Log(ctx, eris.NewFrom(err).BadRequest())
		return
	}

	user, err := repo.QueryByEmail(loginReq.Email)
	if err != nil { //Không tìm thấy user
		logger.Log(ctx, eris.Warning("User not found").UnAuthorized())
		return
	}

	if !pass.CheckPassword(loginReq.Pass, user.Password, "") {
		_, _ = ctx.WriteString("Wrong password")
		return
	}

	//Login thành công thì quay về trang chủ
	_ = ctx.JSON(pmodel.AuthenInfo{
		UserId:       user.Id,
		UserFullName: user.FullName,
		UserEmail:    user.Email,
		Roles:        pmodel.IntArrToRoles(user.Roles), //Chuyển từ mảng []int sang map[int]bool
	})
}

func LogoutREST(ctx iris.Context) {
	_ = session.Logout(ctx)
	_ = ctx.JSON("Logout success")
}
