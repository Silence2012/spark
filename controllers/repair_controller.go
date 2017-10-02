package controllers

import (

	"github.com/astaxie/beego"
	"encoding/json"
	"../constants"
	"errors"
	"../models"
)

type RepairController struct {
	beego.Controller
}


func (this *RepairController) SaveRepairForm() {

	result := make(map[string]interface{})

	body := make(map[string]string)
	bodyJson := this.Ctx.Input.RequestBody

	err := json.Unmarshal(bodyJson, &body)
	this.HandleError(result, err)

	//公司名称
	company, companyExised := body["company"]
	if !companyExised {
		err := errors.New("公司名称必填")
		this.HandleError(result, err)
	}
	//区域，非必填
	region, _ := body["region"]

	//真实姓名（必填）
	name, nameExisted := body["name"]
	if !nameExisted {
		err = errors.New("姓名必填")
		this.HandleError(result, err)
	}
	//手机号码（必填）
	mobile, mobileExisted := body["mobile"]
	if !mobileExisted {
		err = errors.New("手机号必填")
		this.HandleError(result, err)
	}
	//邮箱（必填）
	email, emailExisted := body["email"]
	if !emailExisted {
		err = errors.New("邮箱必填")
		this.HandleError(result, err)
	}
	//行业（必选）
	industry, industryExisted := body["industry"]
	if !industryExisted {
		err = errors.New("行业必填")
		this.HandleError(result, err)
	}
	//产品序列号（必填）
	serial, serialExisted := body["serial"]
	if !serialExisted {
		err = errors.New("产品序列号必填")
		this.HandleError(result, err)
	}
	//设备类型（必选）
	firstDeviceType, firstDeviceTypeExisted := body["firstDeviceType"]
	if !firstDeviceTypeExisted {
		err = errors.New("设备类型必填")
		this.HandleError(result, err)
	}
	secondDeviceType, secondDeviceTypeExisted := body["secondDeviceType"]
	if !secondDeviceTypeExisted {
		err = errors.New("设备类型必填")
		this.HandleError(result, err)
	}
	//寄付帐单地址（必填）
	billAddress, billAddressExisted := body["billAddress"]
	if !billAddressExisted {
		err = errors.New("寄付帐单地址必填")
		this.HandleError(result, err)
	}
	//详细公司地址（必填）
	companyAddress, companyAddressExisted := body["companyAddress"]
	if !companyAddressExisted {
		err = errors.New("详细公司地址必填")
		this.HandleError(result, err)
	}
	//故障细节
	//TODO, 这里要支持语音
	bugDetail, bugDetailExisted := body["bodyDetail"]
	if !bugDetailExisted {
		err = errors.New("故障细节必填")
		this.HandleError(result, err)
	}
	//附件文档
	//TODO, 这里要支持录视频和拍照片， 以及上传文件

	var repairForm map[string]string
	repairForm = make(map[string]string)
	repairForm[constants.Company] = company
	repairForm[constants.Region] = region
	repairForm[constants.Name] = name
	repairForm[constants.Mobile] = mobile
	repairForm[constants.Email] = email
	repairForm[constants.Industry] = industry
	repairForm[constants.Serial] = serial
	repairForm[constants.FirstDeviceType] = firstDeviceType
	repairForm[constants.SecondDeviceType] = secondDeviceType
	repairForm[constants.BillAddress] = billAddress
	repairForm[constants.CompanyAddress] = companyAddress
	repairForm[constants.BugDetail] = bugDetail

	addErr := models.AddRepairForm(repairForm)
	if addErr != nil {
		this.HandleError(result, addErr)
	}

	//生成订单号
	orderPrefix := GenerateOrderPrefix(industry)
	orderNumber, orderErr := models.GenerateRepairFormOrder(orderPrefix)
	if orderErr != nil {
		this.HandleError(result, orderErr)
	}
	this.Ctx.ResponseWriter.Write([]byte(orderNumber))


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



func (this *RepairController) HandleError (result map[string]interface{}, err error) {

	if err != nil {
		this.Ctx.Output.Status = 503
		result[constants.ERROR] = err.Error()
		response, _ := json.Marshal(result)
		this.Ctx.Output.Body(response)
		this.StopRun()
	}

}


