package session

import (
	"context"
	"encoding/json"

	"github.com/TechMaster/core/pmodel"
)

/*
	Theo yêu cầu của Đức, mỗi lần userlogin, lưu thêm một {Key: user.ID, Value: array of sess.ID()} vào Redis
	Chú ý: một người có thể đăng nhập từ nhiều máy tính, thiết bị khác nhau
	Khi hệ thống cập nhật Role của một user từ trang admin
	- Cập nhật vào database Postgresql role của User đó. Phần này không phải trách nhiệm của session
	- Tìm theo key == user.ID để lấy ra mảng các sess.ID()
	  - Xoá tất cả các record có key là sess.ID()
	- Xoá đó xoá nốt record có key == user.ID

	Việc này sẽ buộc người dùng phải đăng nhập lại ở tất cả các thiết bị.
*/

//Hàm này chỉ được chạy bởi Admin
func UpdateRole(userID string, roles []int) error {
	bgCtx := context.Background()
	arrSessID, err := redisClient.SMembers(bgCtx, userID).Result()

	if err != nil {
		return err
	}

	//Cập nhật lại AuthInfo
	for _, sessid := range arrSessID {
		str, err := redisClient.HGet(bgCtx, sessid, SESS_USER).Result()
		if err != nil {
			return err
		}

		var authInfo pmodel.AuthenInfo
		err = json.Unmarshal([]byte(str), &authInfo)
		if err != nil {
			return err
		}
		authInfo.Roles = pmodel.IntArrToRoles(roles)

		var data []byte
		data, err = json.Marshal(authInfo)
		if err != nil {
			return err
		}
		redisClient.HSet(bgCtx, sessid, SESS_USER, string(data))

	}
	return nil
}
