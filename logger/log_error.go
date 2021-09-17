package logger

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/TechMaster/eris"

	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
)

// Chuyên xử lý các err mà controller trả về
func Log(ctx iris.Context, err error) {
	//Trả về JSON error khi client gọi lên bằng AJAX hoặc request.ContentType dạng application/json
	shouldReturnJSON := ctx.IsAjax() || ctx.GetContentTypeRequested() == "application/json"
	switch e := err.(type) {
	case *eris.Error:
		if e.ErrType > eris.WARNING { //Chỉ log ra console hoặc file
			logErisError(e)
		}

		if shouldReturnJSON { //Có trả về báo lỗi dạng JSON cho REST API request không?
			errorBody := iris.Map{
				"error": e.Error(),
			}
			if e.Data != nil { //không có dữ liệu đi kèm thì chỉ cần in thông báo lỗi
				errorBody["data"] = e.Data
			}
			if e.Code > 300 {
				ctx.StatusCode(e.Code)
			} else {
				ctx.StatusCode(iris.StatusInternalServerError)
			}

			_, _ = ctx.JSON(errorBody) //Trả về cho client gọi REST API
			return                     //Xuất ra JSON rồi thì không hiển thị Error Page nữa
		}

		// Nếu request không phải là REST request (AJAX request) thì render error page
		ctx.ViewData("ErrorMsg", e.Error())
		if e.Data != nil {
			if bytes, err := json.Marshal(e.Data); err == nil {
				ctx.ViewData("Data", string(bytes))
			}
		}
		_ = ctx.View(LogConf.ErrorTemplate)
		return		
	default: //Lỗi thông thường
		if errors.Is(err, syscall.EPIPE){
			return
		}
		fmt.Println(err.Error()) //In ra console
		if shouldReturnJSON {    //Trả về JSON
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(err.Error())
		} else {
			_ = ctx.View(LogConf.ErrorTemplate, iris.Map{
				"ErrorMsg": err.Error(),
			})
		}
		return
	}
}

//Kiểm tra xem có lỗi thì báo lỗi
func CheckErr(ctx iris.Context, err error) {
	if err != nil {
		Log(ctx, eris.WrapFrom(err, 4))
	}
}
