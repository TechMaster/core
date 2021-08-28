package template

import (
	"github.com/TechMaster/core/blocks"
	"github.com/TechMaster/core/config"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/view"
)

var HTMLEngine *view.HTMLEngine
var BlockEngine *blocks.BlocksEngine //Đây là
var ViewEngine context.ViewEngine    //generic interface cho các loại view engine

/*
viewFolder: thư mục chứa View Template
defaultLayout: template layout mặc định
*/
func InitHTMLEngine(app *iris.Application, viewFolder string, defaultLayout string) {
	HTMLEngine = iris.HTML(viewFolder, ".html")
	//Nếu app đang debug thì reload bằng true
	HTMLEngine.Layout(defaultLayout).Reload(config.IsAppInDebugMode())
	ViewEngine = HTMLEngine //Gán vào biến này để phần email sẽ dùng
	app.RegisterView(HTMLEngine)
}

func InitBlockEngine(app *iris.Application, viewFolder string, defaultLayout string) {
	BlockEngine = blocks.NewBlocks(viewFolder, ".html")
	//Nếu app đang debug thì reload bằng true
	BlockEngine.Layout(defaultLayout).Reload(config.IsAppInDebugMode())
	ViewEngine = BlockEngine //Gán vào biến này để phần email sẽ dùng
	app.RegisterView(BlockEngine)
}
