package main

import (
	"github.com/astaxie/beego"
	"./routers"
)

func main() {
	routers.InitRepairRouter()
	beego.Run()
}
