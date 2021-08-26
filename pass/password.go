package pass

import (
	"crypto/sha1"

	p "github.com/wuriyanto48/go-pbkdf2"
	"golang.org/x/crypto/bcrypt"
)

/*
Bcrypt băm password rất chậm, tuy nhiên mỗi lần băm, cùng một chuỗi sẽ cho ra kết quả khác nhau
Tránh lỗ hổng bảo mật khi hacker truy theo rainbow table, bảng những password băm sẵn
*/
func HashBcryptPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
		return pass.VerifyPassword(inputpass, hashedpass, hashedpass)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(inputpass))
		return err == nil
	}
}
