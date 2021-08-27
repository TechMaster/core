package logger

import (
	"github.com/kataras/iris/v12"
)

/*
Trang trả về thông tin cho người dùng. Đây không phải là báo lỗi mà chỉ
Nếu redirectLink gô thì hiển thị link để người dùng ấn chuyển tiếp sang trang khác
*/
func Info(ctx iris.Context, msg string, redirectLink ...string) {
	switch len(redirectLink) {
	case 2:
		_ = ctx.View(logConfig.InfoTemplate, iris.Map{
			"Msg":       msg,
			"LinkTitle": redirectLink[1], //<a href='Link'>LinkTitle</a>
			"Link":      redirectLink[0],
		})
	case 1:
		_ = ctx.View(logConfig.InfoTemplate, iris.Map{
			"Msg":       msg,
			"LinkTitle": redirectLink[0], //<a href='Link'>LinkTitle</a>
			"Link":      redirectLink[0],
		})
	default:
		_ = ctx.View(logConfig.InfoTemplate, iris.Map{
			"Msg": msg,
		})
	}
}
