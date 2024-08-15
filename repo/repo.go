package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var Users = make(map[string]*pmodel.User)

func init() {
	CreateNewUser("Admin", "1", "admin@gmail.com", "0123776000", 1)
	CreateNewUser("Bùi Văn Hiên", "1", "hien@gmail.com", "0123456789", 3, 8)
	CreateNewUser("Nguyễn Hàn Duy", "1", "duy@gmail.com", "0123456786", 3, 2)
	CreateNewUser("Phạm Thị Mẫn", "1", "man@gmail.com", "0123456780", 4, 2)
	CreateNewUser("Trịnh Minh Cường", "1", "cuong@gmail.com", "0123456000", 1, 3)
	CreateNewUser("Nguyễn Thành Long", "1", "long@gmail.com", "0123456001", 2)
	CreateNewUser("Dương Văn Thịnh", "1", "thinh@gmail.com", "0223456001", 6, 7, 5, 2)
}

func CreateNewUser(fullName string, password string, email string, phone string, roles ...int) {
	hassedpass := pass.HashPassword(password)
	id, _ := gonanoid.New(8)
	user := pmodel.User{
		Id:       id,
		FullName: fullName,
		Password: hassedpass,
		Email:    strings.ToLower(email),
		Phone:    phone,
		Roles:    roles,
	}

	Users[user.Email] = &user //Thêm user vào users
}
func QueryByEmail(email string) (user *pmodel.User, err error) {
	user = Users[strings.ToLower(email)]
	if user == nil {
		return nil, errors.New("User not found")
	} else {
		return user, nil
	}
}

/*
Có user trong users thì cập nhật: Update
Chưa có thì tạo mới: Insert
Gọi tóm tắt là UpSert
*/
func UpSertUser(user *pmodel.User) {
	Users[strings.ToLower(user.Email)] = user
}

type ViewUser struct {
	Id       string
	FullName string
	Email    string
	Roles    string
}

func GetAll() []ViewUser {
	result := make([]ViewUser, 0, len(Users))
	for _, user := range Users {
		rolesString := ""
		for i, role := range user.Roles {
			rolesString += fmt.Sprintf("%d:%s", role, rbac.RoleName(role))
			if i < len(user.Roles)-1 {
				rolesString += ", "
			}
		}

		viewUser := ViewUser{
			Id:       user.Id,
			FullName: user.FullName,
			Email:    user.Email,
			Roles:    rolesString,
		}

		result = append(result, viewUser)
	}
	return result
}
