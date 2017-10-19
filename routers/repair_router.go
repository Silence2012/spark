package routers

import (
	"github.com/astaxie/beego"
	"../controllers"
)
func InitRepairRouter ()  {
	//新增加报修单
	beego.Router("/repairs/add", &controllers.RepairController{}, "post:SaveRepairForm")
	//按报修单id查询报修单状态
	beego.Router("/repairs/status/:orderid", &controllers.RepairController{}, "get:QueryStatusByOrderId")
	//按报修单id查询报修单详情
	beego.Router("/repairs/query/:orderid", &controllers.RepairController{}, "get:QueryDetailByOrderId")
	//获取所有报修单状态，未处理多少个，正在处理中多少个，已经完成多少个
	beego.Router("/repairs/list/status", &controllers.RepairController{}, "get:GetRepairFormListStatus")
	//按订单状态查询订单列表，未处理new，正在处理handling，已经完成complete
	beego.Router("/repairs/list/:orderstatus", &controllers.RepairController{}, "get:GetRepairFormListByOrderStatus")
	//更新订单状态
	beego.Router("/repairs/update", &controllers.RepairController{}, "post:UpdateRepairForm")
	//查看已完成订单详细情况
	beego.Router("/repairs/complete/detail/:orderid", &controllers.RepairController{}, "get:QueryDetailByOrderId")
	//订单置顶
	beego.Router("/repairs/top", &controllers.RepairController{}, "post:TopOrder")
}