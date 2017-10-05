package models

import (

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"../constants"
	"time"
	"strconv"
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
	//状态 未处理new，正在处理handling，已经完成complete，默认new
	Status string
	//订单号
	OrderId string
	//提交报修单时间
	SubmitTime int64
}

type RepairOrder struct {
	OrderDay string
	OrderNumber int
}

func InitMongodbSession() (*mgo.Session, error) {
	mongoIPs := beego.AppConfig.String("mongodbIPs")
	session, err := mgo.Dial(mongoIPs)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func AddRepairForm(repairFormMap map[string]string) error {
	session, err := InitMongodbSession()
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
			constants.OrderNew,
			repairFormMap[constants.OrderId],
			time.Now().Unix(),
		})
	if err != nil {
		return err
	}
	return nil
}

//生成订单号
func GenerateRepairFormOrder(orderNumberPrefix string) (string, error) {
	year := time.Now().Year()
	month := time.Now().Month()
	monthNumber := int(month)
	day := time.Now().Day()

	yearString := strconv.Itoa(year)
	monthString := strconv.Itoa(monthNumber)
	dayString := strconv.Itoa(day)
	if len(dayString) < 2 {
		dayString = "0" + dayString
	}

	prefix := yearString + monthString + dayString

	session, err := InitMongodbSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairorders")
	_, err = c.Upsert(bson.M{"orderday": prefix}, bson.M{"$inc": bson.M{"ordernumber": 1}})
	if err != nil {
		return "", err
	}

	result := RepairOrder{}
	err = c.Find(bson.M{"orderday": prefix}).One(&result)
	if err != nil {
		return "", err
	}
	tempOrderNumber := strconv.Itoa(result.OrderNumber)
	if len(tempOrderNumber) == 1 {
		tempOrderNumber = "00" + tempOrderNumber
	} else if len(tempOrderNumber) == 2 {
		tempOrderNumber = "0" + tempOrderNumber
	}
	orderNumber := orderNumberPrefix + yearString + monthString + dayString + tempOrderNumber
	return orderNumber, err
}

func GetStatusList() ([]bson.M, error) {
	session, err := InitMongodbSession()
	if err != nil {
		return nil, err
	}

	c := session.DB("ndc").C("repairforms")
	//db.getCollection('repairforms').aggregate({
	//	"$group": {
	//	"_id" : "$status",
	//	"count": {"$sum": 1}
	//	}
	//})

	pipe := c.Pipe([]bson.M{
		{"$group": bson.M{"_id":"$status", "count": bson.M{"$sum": 1}}}})
	resp := []bson.M{}
	pipeErr := pipe.All(&resp)
	beego.Info(resp)
	return resp, pipeErr
}