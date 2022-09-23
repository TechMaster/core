package logger

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12/cache"
	"github.com/kataras/iris/v12/context"
	"syscall"

	"github.com/TechMaster/eris"

	"github.com/goccy/go-json"
	"github.com/kataras/iris/v12"
)

// Chuyên xử lý các err mà controller trả về
func Log(ctx iris.Context, err error) {
	if errors.Is(err, syscall.EPIPE) {
		return
	}
	//Trả về JSON error khi client gọi lên bằng AJAX hoặc request.ContentType dạng application/json
	shouldReturnJSON := ctx.IsAjax() || ctx.GetContentTypeRequested() == "application/json"
	switch e := err.(type) {
	case *eris.Error:
		if e.ErrType > eris.WARNING { //Chỉ log ra console hoặc file
			logErisError(e)
		}

		if shouldReturnJSON { //Có trả về báo lỗi dạng JSON cho REST API request không
			if e.Code > 300 {
				ctx.StatusCode(e.Code)
			} else {
				ctx.StatusCode(iris.StatusInternalServerError)
			}

			_ = ctx.JSON(e.Error()) //Trả về cho client gọi REST API
			return                  //Xuất ra JSON rồi thì không hiển thị Error Page nữa
		}
		setNoCache(ctx)

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
		fmt.Println(err.Error()) //In ra console
		if shouldReturnJSON {    //Trả về JSON
			ctx.StatusCode(iris.StatusInternalServerError)
			_ = ctx.JSON(err.Error())
		} else {
			setNoCache(ctx)
			_ = ctx.View(LogConf.ErrorTemplate, iris.Map{
				"ErrorMsg": err.Error(),
			})
		}
		return
	}
}

func setNoCache(ctx iris.Context) {
	ctx.Header(context.CacheControlHeaderKey, cache.CacheControlHeaderValue)
	ctx.Header(cache.PragmaHeaderKey, cache.PragmaNoCacheHeaderValue)
	ctx.Header(cache.ExpiresHeaderKey, cache.ExpiresNeverHeaderValue)
}
