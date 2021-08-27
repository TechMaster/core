package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TechMaster/core/pass"
	"github.com/TechMaster/core/pmodel"
	"github.com/TechMaster/core/rbac"
	"github.com/segmentio/ksuid"
)

var users = make(map[string]*pmodel.User)

func init() {
	CreateNewUser("Admin", "1", "admin@gmail.com", "0123776000", rbac.ADMIN)
	CreateNewUser("Bùi Văn Hiên", "1", "hien@gmail.com", "0123456789", rbac.TRAINER, rbac.MAINTAINER)
	CreateNewUser("Nguyễn Hàn Duy", "1", "duy@gmail.com", "0123456786", rbac.TRAINER, rbac.STUDENT)
	CreateNewUser("Phạm Thị Mẫn", "1", "man@gmail.com", "0123456780", rbac.SALE, rbac.STUDENT)
	CreateNewUser("Trịnh Minh Cường", "1", "cuong@gmail.com", "0123456000", rbac.ADMIN, rbac.TRAINER)
	CreateNewUser("Nguyễn Thành Long", "1", "long@gmail.com", "0123456001", rbac.STUDENT)
	CreateNewUser("Dương Văn Thịnh", "1", "thinh@gmail.com", "0223456001", rbac.AUTHOR, rbac.EDITOR, rbac.EMPLOYER, rbac.STUDENT)
}

func CreateNewUser(fullName string, password string, email string, phone string, roles ...int) {
	hassedpass, _ := pass.HashBcryptPass(password)

	user := pmodel.User{
		Id:       ksuid.New().String(),
		FullName: fullName,
		Password: hassedpass,
		Email:    strings.ToLower(email),
		Phone:    phone,
		Roles:    roles,
	}

	users[user.Email] = &user //Thêm user vào users
}
func QueryByEmail(email string) (user *pmodel.User, err error) {
	user = users[strings.ToLower(email)]
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
	users[strings.ToLower(user.Email)] = user
}

type ViewUser struct {
	Id       string
	FullName string
	Email    string
	Roles    string
}

func GetAll() []ViewUser {
	result := make([]ViewUser, 0, len(users))
	for _, user := range users {
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
