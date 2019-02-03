package main

import (
	"github.com/astaxie/beego"
	_ "github.com/yacen/song-resgain/web/routers"
	_ "github.com/yacen/song-resgain/web/service"
)

func main() {
	beego.Run()
}
