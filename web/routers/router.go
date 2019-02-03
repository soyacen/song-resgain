package routers

import (
	"github.com/astaxie/beego"
	"github.com/yacen/song-resgain/web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/name", &controllers.NameController{})
}
