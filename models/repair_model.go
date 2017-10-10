package models

import (

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
	"../constants"
	"time"
	"strconv"
)

//报修， 只要客户提交了报修单，那么这里就算完成了
type Report struct {
	//报修时间
	Time   int64 `json:"time"`
	//报修状态
	Complete bool `json:"complete"`
}
//服务中心，这里暂时也这么处理：只要客户提交了报修单，那么这里就算完成了
type Servicecenter struct {
	Time   int64
	Complete bool
}
//工程师上门服务,里面包含了上门服务状态和维修完成状态
type Engineer struct {
	Name            string
	Mobile          string
	//是否已经上门服务
	Homeservice     bool
	//上门服务时间
	Homeservicetime int64

	Notes           string
	//维修是否已完成
	Complete        bool
	//是否发短信通知用户
	Smsuser         bool
	//工程师更新状态时间, 如果时间不为零，说明工程师更新过状态，否则认为工程师没进行任何操作
	Time int64
	//维修时间
	RepairTime int64
}

type OrderLog struct {
	Report *Report
	Servicecenter *Servicecenter
	Engineer *Engineer
}

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
	ThirdDeviceType string
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
	//订单维修完成时间
	FixCompletedTime int64
	//报修单状态
	OrderLog *OrderLog
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
			repairFormMap[constants.ThirdDeviceType],
			repairFormMap[constants.BillAddress],
			repairFormMap[constants.CompanyAddress],
			repairFormMap[constants.BugDetail],
			constants.OrderNew,
			repairFormMap[constants.OrderId],
			time.Now().Unix(),
			0,
			&OrderLog{
				&Report{time.Now().Unix(), true},
			    &Servicecenter{time.Now().Unix(), true},
			    &Engineer{},
			},
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

func GetRepairOrderStatus(orderId string) (interface{}, error) {
	session, err := InitMongodbSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairforms")

	var result interface{}
	err = c.Find(bson.M{"orderid": orderId}).Select(bson.M{"orderlog": 1,"status":1,"fixcompletedtime":1}).One(&result)
	if err != nil {
		return nil, err
	}
	beego.Info(result)
	return &result, nil
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

func UpdateOrderLog(body map[string]string) error {
	orderId := body[constants.OrderId]
	engineerName := body[constants.EngineerName]
	engineerMobile := body[constants.EngineerMobile]
	//这个需要转换为int64
	homeServiceTime := body[constants.HomeServiceTime]
	notes := body[constants.Notes]
	//这个需要转换为bool
	fixCompleted := body[constants.FixCompleted]
	//这个需要转换为bool
	smsUser := body[constants.SMSUser]
	//维修时间
	repairTime := body[constants.RepairTime]

	var homeServiceTimeInt int64
	var fixCompletedBool bool
	var smsUserBool bool

	var homeService bool

	fixCompletedBool, _ = strconv.ParseBool(fixCompleted)
	smsUserBool, _ = strconv.ParseBool(smsUser)

	var fixCompletedTime int64
	fixCompletedTime = 0

	var repairTimeInt int64
	repairTimeInt = 0

	if homeServiceTime == "" {
		homeServiceTimeInt = 0
	} else {
		//转化所需模板
		timeLayout := "2006-01-02"
		//获取时区
		loc, _ := time.LoadLocation("Local")
		theTime, convertErr := time.ParseInLocation(timeLayout, homeServiceTime, loc)
		if convertErr != nil {
			homeServiceTimeInt = 0
		}
		//转化为时间戳 类型是int64
		homeServiceTimeInt = theTime.Unix()
	}
	if repairTime == "" {
		repairTimeInt = 0
	} else {
		//转化所需模板
		timeLayout := "2006-01-02"
		//获取时区
		loc, _ := time.LoadLocation("Local")
		theTime, convertErr := time.ParseInLocation(timeLayout, repairTime, loc)
		if convertErr != nil {
			repairTimeInt = 0
		}
		//转化为时间戳 类型是int64
		repairTimeInt = theTime.Unix()
	}

	if homeServiceTimeInt == 0 {
		homeService = false
	}else {
		homeService = true
	}

	var handleStatus string
	if fixCompletedBool {
		handleStatus = constants.OrderCompleted
		fixCompletedTime = time.Now().Unix()
	}else {
		handleStatus = constants.OrderHandling
	}

	session, err := InitMongodbSession()
	if err != nil {
		return  err
	}
	defer session.Close()

	c := session.DB("ndc").C("repairforms")

	updateErr := c.Update(bson.M{"orderid": orderId}, bson.M{"$set": bson.M{
														"orderlog.engineer.name": engineerName,
														"orderlog.engineer.mobile": engineerMobile,
														"orderlog.engineer.homeservice": homeService,
														"orderlog.engineer.homeservicetime": homeServiceTimeInt,
														"orderlog.engineer.notes": notes,
														"orderlog.engineer.complete": fixCompletedBool,
														"orderlog.engineer.smsuser": smsUserBool,
														"orderlog.engineer.time": time.Now().Unix(),
														"orderlog.engineer.repairtime": repairTimeInt,
														"fixcompletedtime": fixCompletedTime,
														"status": handleStatus,
														}})

	return updateErr

}