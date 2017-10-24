package utils

import (
	"net/http"
	"github.com/astaxie/beego"
	"encoding/json"
	"bytes"
	"io/ioutil"
)
func SendSms(url string,PhoneNum string,OrderID string) string {
	beego.Info("URL:>", url , "phoneNums:>", PhoneNum, "orderID:>", OrderID)
	sendContext := map[string]string{"phoneNums": PhoneNum, "orderID": OrderID}

	jsonStr, _ := json.Marshal(sendContext)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error(err)
	}
	defer resp.Body.Close()
	beego.Info("response Status:", resp.Status)
	beego.Info("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
