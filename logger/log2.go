package logger

import (
	"fmt"

	"github.com/TechMaster/eris"
)

/* Log lỗi mà không trả về cho web client. Chỉ dùng Log2 trong những trường hợp
hàm không nằm trong chuỗi các hàm xử lý HTTP request, không có iris.Context
*/
func Log2(err error) {
	switch e := err.(type) {
	case *eris.Error: //Lỗi kiểu eris
		logErisError(e)
	default: //Lỗi thông thường
		fmt.Println(err.Error()) //In ra console
	}
}
