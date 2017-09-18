package main

import (
	"github.com/astaxie/beego"
	"./controllers"
)

func main() {
	beego.Router("/feedback", &controllers.FeedbackController{}, "post:SaveFeedbackForm")
	beego.Run()
}
