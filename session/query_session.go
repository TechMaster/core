package session

import (
	"github.com/TechMaster/core/pmodel"

	"github.com/TechMaster/core/sessions"
	"github.com/TechMaster/eris"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

func GetAuthInfo(ctx iris.Context) (*pmodel.AuthenInfo, error) {
	data := sessions.Get(ctx).Get(SESS_USER)
	if data == nil {
		return nil, nil
	}

	authinfo := new(pmodel.AuthenInfo)
	if err := mapstructure.WeakDecode(data, authinfo); err != nil {
		return nil, eris.NewFrom(err)
	}
	return authinfo, nil
}

/*
Lấy AuthInfo từ trong ViewData.
Ở Handler RBAC checkpermission trước, nếu người dùng đăng nhập
AuthInfo đã được ghi vào ViewData[AUTHINFO]
*/
func GetAuthInfoViewData(ctx iris.Context) *pmodel.AuthenInfo {
	if raw_authinfo := ctx.GetViewData()[AUTHINFO]; raw_authinfo != nil {
		if authinfo, ok := raw_authinfo.(*pmodel.AuthenInfo); ok {
			return authinfo
		}
	}
	return nil
}

func IsLogin(ctx iris.Context) bool {
	login, _ := sessions.Get(ctx).GetBoolean(SESS_AUTH)
	return login
}
