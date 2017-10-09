package controllers

import (

	"github.com/astaxie/beego"
	"encoding/json"
	"../constants"
	"errors"
	"../models"
	"../utils"

	"strings"
)

type RepairController struct {
	beego.Controller
}

type RepairFormStatusList []struct {
	ID    string `json:"_id"`
	Count int    `json:"count"`
}

var titleArray = []string{"公司名称", "区域", "真实姓名", "手机号码", "邮箱", "行业", "产品序列号", "设备类型", "寄付帐单地址", "详细公司地址", "故障细节"}

func (this *RepairController) SaveRepairForm() {

	result := make(map[string]interface{})

	body := make(map[string]string)
	bodyJson := this.Ctx.Input.RequestBody
	beego.Info(string(bodyJson))

	err := json.Unmarshal(bodyJson, &body)
	this.HandleError(result, err)
	//验证输入项
	requestDataArray, validErr := validRepairForm(body)
	this.HandleError(result, validErr)

	//industry, _ := body["industry"]
	firstDeviceType, _ := body["firstDeviceType"]

	//生成订单号前缀
	orderPrefix := GenerateOrderPrefix(firstDeviceType)

	//生成订单号
	orderNumber, orderErr := models.GenerateRepairFormOrder(orderPrefix)
	if orderErr != nil {
		this.HandleError(result, orderErr)
	}
	body[constants.OrderId] = orderNumber
	//持久化输入项
	addErr := models.AddRepairForm(body)
	if addErr != nil {
		this.HandleError(result, addErr)
	}

	//生成excel,文件名就是订单号，保存到本地
	excelPath := generateExcel(requestDataArray, orderNumber)
	//发送邮件
	sendEmail(requestDataArray, excelPath,orderNumber)
	//发送短信

	this.Ctx.ResponseWriter.Write([]byte(orderNumber))

}

//按报修单id查询报修单状态
func (this *RepairController) QueryStatusByOrderId()  {
	result := make(map[string]interface{})

	paramMap := this.Ctx.Input.Params()
	orderId, orderIdExisted := paramMap[":orderid"]
	if !orderIdExisted {
		err := errors.New("请输入订单号")
		this.HandleError(result, err)
	}
	status, getStatusErr := models.GetRepairOrderStatus(orderId)
	this.HandleError(result, getStatusErr)

	response, marshalErr := json.Marshal(status)
	this.HandleError(result, marshalErr)
	this.Ctx.ResponseWriter.Write(response)
}

//按报修单id查询报修单详情
func (this *RepairController) QueryDetailByOrderId()  {
	result := make(map[string]interface{})

	paramMap := this.Ctx.Input.Params()
	orderId, orderIdExisted := paramMap[":orderid"]
	if !orderIdExisted {
		err := errors.New("请输入订单号")
		this.HandleError(result, err)
	}
	detail, detailErr := models.GetRepairOrderDetail(orderId)
	this.HandleError(result, detailErr)

	response, marshalErr := json.Marshal(detail)
	this.HandleError(result, marshalErr)
	this.Ctx.ResponseWriter.Write(response)
}
//获取所有报修单状态，未处理多少个，正在处理中多少个，已经完成多少个
func (this *RepairController) GetRepairFormListStatus()  {

	//db.getCollection('repairforms').aggregate({
	//	"$group": {
	//	"_id" : "$status",
	//	"count": {"$sum": 1}
	//	}
	//})

	result := make(map[string]interface{})
	resp, getStatuListErr := models.GetStatusList()
	this.HandleError(result, getStatuListErr)

	//[map[_id:completed count:4] map[_id:new count:48] map[_id:handling count:3]]
	for _, mapData := range resp {
		id := mapData["_id"]
		count := mapData["count"]
		result[id.(string)] = count
	}
	response, marshalErr := json.Marshal(result)
	this.HandleError(result, marshalErr)
	this.Ctx.ResponseWriter.Write(response)

}
//按订单状态查询订单列表，未处理new，正在处理handling，已经完成complete
func (this *RepairController) GetRepairFormListByOrderStatus()  {
	result := make(map[string]interface{})

	paramMap := this.Ctx.Input.Params()
	orderStatus, orderStatusExisted := paramMap[":orderstatus"]
	if !orderStatusExisted {
		err := errors.New("请输入订单状态")
		this.HandleError(result, err)
	}

	orderIds, getOrderStatusErr := models.GetRepairFormListByOrderStatus(orderStatus)
	this.HandleError(result, getOrderStatusErr)

	response, marshalErr := json.Marshal(orderIds)
	this.HandleError(result, marshalErr)
	this.Ctx.ResponseWriter.Write(response)
}

//更新订单状态
func (this *RepairController) UpdateRepairForm()  {
	result := make(map[string]interface{})

	body := make(map[string]string)
	bodyJson := this.Ctx.Input.RequestBody
	beego.Info(string(bodyJson))
	err := json.Unmarshal(bodyJson, &body)
	this.HandleError(result, err)

	validErr := validEngineerOperations(body)
	this.HandleError(result, validErr)

	//更新订单状态

	getOrderStatusErr := models.UpdateOrderLog(body)
	this.HandleError(result, getOrderStatusErr)

	this.Ctx.ResponseWriter.Write([]byte(constants.SUCCESS))

}


func validRepairForm(body map[string]string) ([]string, error)  {

	//公司名称
	company, companyExised := body[constants.Company]
	if !companyExised {
		err := errors.New("公司名称必填")
		return nil, err
	}
	//区域，非必填
	region, _ := body[constants.Region]

	//真实姓名（必填）
	name, nameExisted := body[constants.Name]
	if !nameExisted {
		err := errors.New("姓名必填")
		return nil, err
	}
	//手机号码（必填）
	mobile, mobileExisted := body[constants.Mobile]
	if !mobileExisted {
		err := errors.New("手机号必填")
		return nil, err
	}
	//邮箱（必填）
	email, emailExisted := body[constants.Email]
	if !emailExisted {
		err := errors.New("邮箱必填")
		return nil, err
	}
	//行业（必选）
	industry, industryExisted := body[constants.Industry]
	if !industryExisted {
		err := errors.New("行业必填")
		return nil, err
	}
	//产品序列号（必填）
	serial, serialExisted := body[constants.Serial]
	if !serialExisted {
		err := errors.New("产品序列号必填")
		return nil, err
	}
	//设备类型（必选）
	firstDeviceType, firstDeviceTypeExisted := body[constants.FirstDeviceType]
	if !firstDeviceTypeExisted {
		err := errors.New("设备类型必填")
		return nil, err
	}
	secondDeviceType, secondDeviceTypeExisted := body[constants.SecondDeviceType]
	if !secondDeviceTypeExisted {
		err := errors.New("设备类型必填")
		return nil, err
	}
	thirdDeviceType, _ := body[constants.ThirdDeviceType]

	//寄付帐单地址（必填）
	billAddress, billAddressExisted := body[constants.BillAddress]
	if !billAddressExisted {
		err := errors.New("寄付帐单地址必填")
		return nil, err
	}
	//详细公司地址（必填）
	companyAddress, companyAddressExisted := body[constants.CompanyAddress]
	if !companyAddressExisted {
		err := errors.New("详细公司地址必填")
		return nil, err
	}
	//故障细节
	//TODO, 这里要支持语音
	bugDetail, bugDetailExisted := body["bugDetail"]
	if !bugDetailExisted {
		err := errors.New("故障细节必填")
		return nil, err
	}
	//附件文档
	//TODO, 这里要支持录视频和拍照片， 以及上传文件

	var result []string
	result = make([]string, 11)
	result[0] = company
	result[1] = region
	result[2] = name
	result[3] = mobile
	result[4] = email
	result[5] = industry
	result[6] = serial
	result[7] = firstDeviceType + secondDeviceType + thirdDeviceType
	result[8] = billAddress
	result[9] = companyAddress
	result[10] = bugDetail
	return result, nil
}

func GenerateOrderPrefix( catalog string) string {
	//化工品和药品	订单生成首字母T
	//涂布复合	订单生成首字母P
	//薄膜和片材挤出	订单生成首字母P
	//食品加工	订单生成首字母T
	//冶金工业	订单生成首字母M
	//矿石和松散物	订单生成首字母T
	//无纺布和纺织品	订单生成首字母P
	//软管及管材	订单生成首字母C
	//橡胶和乙烯基压延	订单生成首字母P
	//烟草加工	订单生成首字母T
	//电线，电缆和光纤	订单生成首字母C
	//其他	订单生成首字母O
//new needs
//	冶金工业系列产品——M；
//	系统选项——S；
//	红外传感器系列产品——I；
//	Beta LaserMike和Zmike 系列产品——C；

/*
	switch industry {
	case "化工品和药品":
		return "T"
	case "涂布复合":
		return "P"
	case "薄膜和片材挤出":
		return "P"
	case "食品加工":
		return "T"
	case "冶金工业":
		return "M"
	case "矿石和松散物":
		return "T"
	case "无纺布和纺织品":
		return "P"
	case "软管及管材":
		return "C"
	case "橡胶和乙烯基压延":
		return "P"
	case "烟草加工":
		return "T"
	case "电线，电缆和光纤":
		return "C"
	default:
		return "O"
	}
	*/
	switch catalog {
	case "冶金工业系列产品":
		return "M"
	case "系统选项":
		return "S"
	case "红外传感器系列产品":
                return "I"
	case "Beta LaserMike和Zmike 系列产品":
		return "C"
	default:
		return "O"
	}
}

func generateExcel(body []string, orderId string) string {

	var content map[int][]string
	content = make(map[int][]string)
	content[0] = body

	var excelData []map[int][]string
	excelData = make([]map[int][]string, 1)
	excelData[0] = content

	excelPath := beego.AppConfig.String(constants.ExcelDir) + orderId + ".xlsx"
	utils.GenerateExcel(titleArray, excelData, excelPath)

	return excelPath
}

func sendEmail(requestDataArray []string, excelPath string, orderNumber string )  {
	//from string, to []string, cc string, subject string, contentType string, body string, attachments ...string

	from := beego.AppConfig.String(constants.EmailUser)
	//发送邮件给相关人员，邮箱配置在配置文件里以逗号隔开
	toStringWithComma := beego.AppConfig.String(constants.EmailList)
	to := strings.Split(toStringWithComma, ",")
	cc := ""
	//TODO 邮件标题,这个也需要最终定了以后替换
	subject := "新的报修单(来自微信、网页端的测试通知邮件）【测试】" + "订单号： " + orderNumber
	contentType := "text/html"

	tbody := ""
	for index, requestData := range requestDataArray {
		title := titleArray[index]
		content := requestData
		tbody = tbody + `<tr style="height: 45px;">
							<td style="height: 45px; width: 312.367px;">
								<pre style="background-color: #ffffff; color: #000000; font-family: 'Menlo'; font-size: 10.5pt;"><span style="color: #008000; font-weight: bold;">`+title+`</span></pre>
							</td>
							<td style="height: 45px; width: 871.633px;">`+content+`</td>
						</tr>`
	}

	body := `<table style="width: 1201px; height: 612px;">
				<tbody>
					`+tbody+`
				</tbody>
			</table>
			<p>&nbsp;</p>`

	attachment := []string{excelPath}

	utils.SendEmail(from, to, cc, subject, contentType, body, attachment...)
}


//处理报修单应该只有报修单号，工程师姓名和电话是必填项，其他可以不填
func validEngineerOperations(body map[string]string) error  {

	//报修单号
	_, orderIdExised := body[constants.OrderId]
	if !orderIdExised {
		err := errors.New("订单号必填")
		return err
	}
	//工程师姓名
	_, engineerNameExisted := body[constants.EngineerName]
	if !engineerNameExisted {
		err := errors.New("工程师姓名必填")
		return err
	}

	//工程师电话
	_, engineerMobileExisted := body[constants.EngineerMobile]
	if !engineerMobileExisted {
		err := errors.New("工程师电话必填")
		return err
	}

	return nil
}


func (this *RepairController) HandleError (result map[string]interface{}, err error) {

	if err != nil {
		this.Ctx.Output.Status = 503
		result[constants.ERROR] = err.Error()
		response, _ := json.Marshal(result)
		this.Ctx.Output.Body(response)
		this.StopRun()
	}

}


