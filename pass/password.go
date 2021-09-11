package pass

import (
	"crypto/sha1"

	p "github.com/wuriyanto48/go-pbkdf2"
)

/*
Tạo interface cho các thư viện băm password tuân thủ
*/
type PasswordLib interface {
	Hash(password string) (hash string)
	Compare(password string, hashpass string) bool
}

var PassLib PasswordLib

func init() {
	PassLib = Argon2id{ //Chọn Argon2id làm thư viện băm password, trong tương lai có thể thay đổi
		params: Params{
			Memory:      64 * 1024,
			Iterations:  1,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

//Sử dụng thuật toán băm Argon2Id
func HashPassword(inputpass string) string {
	return PassLib.Hash(inputpass)
}

/*
Hàm check password hỗ trợ cả kiểu SHA1 cũ và bcrypt mới
- inputpass: password nhập vào lúc login
- hashedpass: password đã băm đọc từ CSDL
- salt: chuỗi nhiễu tạo ra từ thuật toán SHA1 cũ
*/
func CheckPassword(inputpass string, hashedpass string, salt string) bool {
	if salt != "" {
		pass := p.NewPassword(sha1.New, 50, 64, 10000)
		return pass.VerifyPassword(inputpass, hashedpass, salt) //Sửa theo yêu cầu Nhật Đức
	} else {
		return PassLib.Compare(inputpass, hashedpass)
	}
}
