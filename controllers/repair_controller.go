package controllers

import (

	"github.com/astaxie/beego"
	"encoding/json"
	"../constants"
	"errors"
	"../models"
	"../utils"
)

type RepairController struct {
	beego.Controller
}

var titleArray = []string{"公司名称", "区域", "真实姓名", "手机号码", "邮箱", "行业", "产品序列号", "设备类型", "寄付帐单地址", "详细公司地址", "故障细节"}

func (this *RepairController) SaveRepairForm() {

	result := make(map[string]interface{})

	body := make(map[string]string)
	bodyJson := this.Ctx.Input.RequestBody

	err := json.Unmarshal(bodyJson, &body)
	this.HandleError(result, err)
	//验证输入项
	requestDataArray, validErr := validRepairForm(body)
	this.HandleError(result, validErr)

	industry, _ := body["industry"]
	//生成订单号前缀
	orderPrefix := GenerateOrderPrefix(industry)

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
	sendEmail(requestDataArray, excelPath)
	//发送短信

	this.Ctx.ResponseWriter.Write([]byte(orderNumber))

}

func (this *RepairController) QueryStatusByOrderId()  {

}

func (this *RepairController) QueryDetailByOrderId()  {

}

func (this *RepairController) GetRepairFormListStatus()  {

	//db.getCollection('repairforms').aggregate({
	//	"$group": {
	//	"_id" : "$status",
	//	"count": {"$sum": 1}
	//	}
	//})

}

func (this *RepairController) GetRepairFormListByOrderStatus()  {

}

func (this *RepairController) UpdateRepairForm()  {

}

func (this *RepairController) QueryCompletedRepairFormDetailByOrderId()  {

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

	result := []string{}
	result[0] = company
	result[1] = region
	result[2] = name
	result[3] = mobile
	result[4] = email
	result[5] = industry
	result[6] = serial
	result[7] = firstDeviceType + secondDeviceType
	result[8] = billAddress
	result[9] = companyAddress
	result[10] = bugDetail
	return result, nil
}

func GenerateOrderPrefix(industry string) string {
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

func sendEmail(requestDataArray []string, excelPath string)  {
	//from string, to []string, cc string, subject string, contentType string, body string, attachments ...string

	from := beego.AppConfig.String(constants.EmailUser)
	//TODO 还没拿到发送的邮箱，暂时写个测试数据
	to := []string{"joey8656@163.com"}
	cc := ""
	//TODO 邮件标题,这个也需要最终定了以后替换
	subject := "新的报修单"
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


func (this *RepairController) HandleError (result map[string]interface{}, err error) {

	if err != nil {
		this.Ctx.Output.Status = 503
		result[constants.ERROR] = err.Error()
		response, _ := json.Marshal(result)
		this.Ctx.Output.Body(response)
		this.StopRun()
	}

}


