
package utils


import (
	"github.com/go-redis/redis"
	"../constants"
	"github.com/astaxie/beego"
	"time"
)



var client *redis.Client

const (
	TokenKey = "WeixinToken"
	TicketKey = "WeixinTicket"
	ExpireTime = 60 * 60 * 1.5 * time.Second
)

func initClient() {
	var ip = ""
	var port = ""
	var password = ""

	ip = beego.AppConfig.String(constants.RedisIp)
	port = beego.AppConfig.String(constants.RedisPort)
	password = beego.AppConfig.String(constants.RedisPassword)

	baseConfig := &redis.Options{
		Addr:         ip + ":"+ port,
		Password:     password,
	}
	client = redis.NewClient(baseConfig)

}


func WriteAccessToken(token string) {
	if client == nil {
		initClient()
	}
	client.Set(TokenKey, token, ExpireTime)

}

func GetAccessToken() (string, error) {
	if client == nil {
		initClient()
	}
	result := client.Get(TokenKey)
	accessToken, accessTokenErr := result.Result()
	return accessToken, accessTokenErr
}

func WriteTicket(ticket string) {
	if client == nil {
		initClient()
	}
	client.Set(TicketKey, ticket, ExpireTime)

}

func GetTicket() (string, error) {
	if client == nil {
		initClient()
	}
	result := client.Get(TicketKey)
	ticket, ticketErr := result.Result()
	return ticket, ticketErr
}
