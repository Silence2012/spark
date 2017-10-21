package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
)

type UserInfo struct {
	OpenId string `json:"openid"`
	NickName string `json: "nickname"`
	Sex int `json:"sex"`
	Language string `json:"language"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	Privilege []string `json: "privilege"`
}

func AddWeixinUserInfo(userData UserInfo) error {
	beego.Info("add weixin user info....")
	beego.Info(userData)
	session, err := InitMongodbSession()
	if err != nil {
		return  err
	}
	defer session.Close()

	c := session.DB("ndc").C("weixin_user")

	_, updateErr := c.Upsert(bson.M{"openid": userData.OpenId}, bson.M{"$set": bson.M{
		"openid": userData.OpenId,
		"nickname": userData.NickName,
		"sex": userData.Sex,
		"language": userData.Language,
		"city": userData.City,
		"province": userData.Province,
		"country": userData.Country,
		"headimgurl": userData.HeadImgUrl,
		"privilege": userData.Privilege,
	}})

	if updateErr != nil {
		return updateErr
	}
	return nil
}
