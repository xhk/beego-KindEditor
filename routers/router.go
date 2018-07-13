package routers

import (
	"beego-KindEditor/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/uploadfilemgr", &controllers.UploadFileMgrController{})
}
