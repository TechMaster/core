package session

import (
	"context"

	"github.com/TechMaster/core/pmodel"

	"github.com/TechMaster/core/sessions"
	"github.com/kataras/iris/v12"
)

/*
Thực hiện sau khi người dùng login thành tạo trong session
những key/value chứa thông tin người đăng nhập
*/
func SetAuthenticated(ctx iris.Context, authenInfo pmodel.AuthenInfo) error {
	sess := sessions.Get(ctx)
	sess.Set(SESS_USER, authenInfo)

	/* Thêm sessionID vào key userID
	Ứng với userID sẽ có một Set các sessionID
	*/
	bgCtx := context.Background()
	_, err := redisClient.SAdd(bgCtx, authenInfo.Id, sess.ID()).Result()
	if err != nil {
		return err
	}
	redisClient.Expire(bgCtx, authenInfo.Id, expires) //Đặt thời điểm hết hạn cho bản ghi này
	return nil
}

/*
Trong Redis ứng với mỗi user đăng nhập trên thiết bị A và B có những record sau:
SessionID trên A : userID,...
SessionID trên B : userID,...
userID: SessionID A, SessionID B

Khi user logout ở trên thiết bị A thì cần xoá
SessionID trên A : userID,...
và cập nhật lại userID: SessionID B

Như vậy Redis không bị rác
*/
func Logout(ctx iris.Context) error {
	authenInfo := GetAuthInfo(ctx)
	sess := sessions.Get(ctx)
	bgCtx := context.Background()
	sessionID := sess.ID()
	redisClient.Del(bgCtx, sessionID) //Xoá sessionID thực sự

	if authenInfo != nil {
		//Loại bớt một phần tử trong tập một user.Id chứa nhiều session id của một user
		_, err := redisClient.SRem(bgCtx, authenInfo.Id, sessionID).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
