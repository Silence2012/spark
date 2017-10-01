package models

import (
	"log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"../constants"
)

type RepairForm struct {
	//公司名称
	Company string
	//区域，非必填
	Region string
	//真实姓名（必填）
	Name string
	//手机号码（必填）
	Mobile string
	//邮箱（必填）
	Email string
	//行业（必选）
	Industry string
	//产品序列号（必填）
	Serial string
	//设备类型（必选）
	FirstDeviceType string
	SecondDeviceType string
	//寄付帐单地址（必填）
	BillAddress string
	//详细公司地址（必填）
	CompanyAddress string
	//故障细节
	//TODO, 这里要支持语音
	BugDetail string
	//附件文档
	//TODO, 这里要支持录视频和拍照片， 以及上传文件
}

func AddRepairForm(repairFormMap map[string]string) error {
	mongoIPs := beego.AppConfig.String("mongodbIPs")
	session, err := mgo.Dial(mongoIPs)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairforms")
	err = c.Insert(&RepairForm{
			repairFormMap[constants.Company],
			repairFormMap[constants.Region],
			repairFormMap[constants.Name],
			repairFormMap[constants.Mobile],
			repairFormMap[constants.Email],
			repairFormMap[constants.Industry],
			repairFormMap[constants.Serial],
			repairFormMap[constants.FirstDeviceType],
			repairFormMap[constants.SecondDeviceType],
			repairFormMap[constants.BillAddress],
			repairFormMap[constants.CompanyAddress],
			repairFormMap[constants.BugDetail],
		})
	if err != nil {
		return err
	}
	return nil
}