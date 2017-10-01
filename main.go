package main

import (
	"github.com/astaxie/beego"
	"./controllers"
)

func main() {
	beego.Router("/repairs", &controllers.RepairController{}, "post:SaveRepairForm")
	beego.Run()
}
