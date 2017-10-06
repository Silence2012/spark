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
	//报修单状态
	OrderLog         struct {
		//报修， 只要客户提交了报修单，那么这里就算完成了
		Report struct {
			//报修时间
			Time   int64 `json:"time"`
			//报修状态
			Complete bool `json:"complete"`
		} `json:"report"`
		//服务中心，这里暂时也这么处理：只要客户提交了报修单，那么这里就算完成了
		Servicecenter struct {
			Time   string `json:"time"`
			Complete bool `json:"complete"`
		} `json:"servicecenter"`
		//工程师上门服务,里面包含了上门服务状态和维修完成状态
		Engineer struct {
			Name            string `json:"name"`
			Mobile          string `json:"mobile"`
			//是否已经上门服务
			Homeservice     bool   `json:"homeservice"`
			//上门服务时间
			Homeservicetime int64 `json:"homeservicetime"`
			Notes           string `json:"notes"`
			//维修是否已完成
			Complete        bool   `json:"complete"`
			//是否发短信通知用户
			Smsuser         bool   `json:"smsuser"`
			//工程师更新状态时间
			Time int64 `json:"time"`
		} `json:"engineer"`
	} `json:"orderLog"`
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
	defer session.Close()

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

func GetRepairOrderDetail(orderId string) (*RepairForm, error) {
	session, err := InitMongodbSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairforms")
	result := RepairForm{}
	err = c.Find(bson.M{"orderid": orderId}).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetRepairFormListByOrderStatus(status string) ([]interface{}, error)  {
	session, err := InitMongodbSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairforms")
	var result []interface{}
	err = c.Find(bson.M{"status": status}).Select(bson.M{"orderid": 1}).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}