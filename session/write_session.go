package session

import (
	"github.com/TechMaster/core/pmodel"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
)

/*
Thực hiện sau khi người dùng login thành tạo trong session
những key/value chứa thông tin người đăng nhập
*/
func SetAuthenticated(ctx iris.Context, authenInfo pmodel.AuthenInfo) {
	sess := sessions.Get(ctx)
	sess.Set(SESS_AUTH, true)
	sess.Set(SESS_USER, authenInfo)
}
