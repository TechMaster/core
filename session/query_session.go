package session

import (
	"github.com/TechMaster/core/pmodel"

	"github.com/TechMaster/core/sessions"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

/*
Lấy thông tin đăng nhập của người dùng từ Redis Session
*/
func GetAuthInfoSession(ctx iris.Context) (authinfo *pmodel.AuthenInfo) {
	data := sessions.Get(ctx).Get(SESS_USER)
	if data == nil {
		return nil
	}

	authinfo = new(pmodel.AuthenInfo)
	if err := mapstructure.WeakDecode(data, authinfo); err != nil {
		return nil
	}
	return authinfo
}

/*
Lấy AuthInfo từ trong ViewData. Nếu lấy được thì trả về luôn.
Nếu không tồn tại kiểm tra tiếp trong session. Nếu thấy trả về.
Nếu trong ViewData và Session đều không có, có nghĩa người dùng chưa login
*/
func GetAuthInfo(ctx iris.Context) (authinfo *pmodel.AuthenInfo) {
	if raw_authinfo := ctx.GetViewData()[AUTHINFO]; raw_authinfo != nil {
		var ok bool
		if authinfo, ok = raw_authinfo.(*pmodel.AuthenInfo); ok {
			return authinfo
		}
	}

	authinfo = GetAuthInfoSession(ctx)
	if authinfo != nil {
		return authinfo
	}
	return nil
}
