package template

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/view"
)

var ViewEngine *view.BlocksEngine

//Nhật Đức chỉnh sửa, thay đổi view template tùy chỉnh theo site thay vì cố định như trước.
func InitViewEngine(app *iris.Application, view string) {
	ViewEngine = iris.Blocks(view, ".html")
	ViewEngine.Layout("layout/layout.html").Reload(true)	
	app.RegisterView(ViewEngine)
}
