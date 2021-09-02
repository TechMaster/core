package controller

import (
	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/TechMaster/core/repo"
	"github.com/TechMaster/eris"

	"github.com/TechMaster/core/session"

	"github.com/kataras/iris/v12"
)

func ShowHomePage(ctx iris.Context) {
	authinfo := session.GetAuthInfo(ctx)
	if authinfo != nil {
		ctx.ViewData("roles", rbac.RolesNames(authinfo.Roles))
	}

	ctx.ViewData("users", repo.GetAll())
	_ = ctx.View("index")
}

func ShowSecret(ctx iris.Context) {
	logger.Info(ctx, "Đây là trang bí mật chỉ dành cho người đã đăng nhập")
}

func ShowErr(ctx iris.Context) {
	if err := foo(); err != nil {
		logger.Log(ctx, err)
		return
	}
	logger.Info(ctx, "Không có lỗi gì cả")
}

func foo() error {
	if err := bar(); err != nil {
		return err
	}
	return nil
}

func bar() error {
	return eris.SysError("Show Stack Error")
}

func Books(ctx iris.Context) {
	authinfo := session.GetAuthInfo(ctx)
	if authinfo == nil {
		logger.Log(ctx, eris.Warning("Bạn chưa đăng nhập").UnAuthorized())
	}

	type Book struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	type Data struct {
		AuthInfo *pmodel.AuthenInfo `json:"authinfo"`
		Books    []Book             `json:"books"`
	}

	_, _ = ctx.JSON(Data{
		AuthInfo: authinfo,
		Books: []Book{
			{
				Title:  "Dế Mèn Phiêu Lưu Ký",
				Author: "Tô Hoài",
			},
			{
				Title:  "Nhật Ký Trong Tù",
				Author: "Hồ Chí Minh",
			},
			{
				Title:  "Tắt Đèn",
				Author: "Ngô Tất Tố",
			},
		},
	})
}
