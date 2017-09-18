package controllers

import (

	"github.com/astaxie/beego"
	"encoding/json"
	"../constants"
	"errors"
)

type FeedbackController struct {
	beego.Controller
}


func (this *FeedbackController) SaveFeedbackForm() {
	//id := this.Ctx.Input.Param(":id")
	result := make(map[string]interface{})

	body := make(map[string]string)
	bodyJson := this.Ctx.Input.RequestBody

	err := json.Unmarshal(bodyJson, &body)
	this.HandleError(result, err)

	name, nameExised := body["name"]
	if !nameExised {
		err := errors.New("姓名必填")
		this.HandleError(result, err)
	}

	result[constants.REUSLT] = constants.SUCCESS
	result[constants.DETAIL] = name
	response, _ := json.Marshal(result)
	this.Ctx.ResponseWriter.Write(response)

	//var t models.Task
	//if err := json.Unmarshal(this.Ctx.Input.RequestBody, &t); err != nil {
	//	this.Ctx.Output.SetStatus(400)
	//	this.Ctx.Output.Body([]byte(err.Error()))
	//	return
	//}
	//if t.ID != intid {
	//	this.Ctx.Output.SetStatus(400)
	//	this.Ctx.Output.Body([]byte("inconsistent task IDs"))
	//	return
	//}
}



func (this *FeedbackController) HandleError (result map[string]interface{}, err error) {

	if err != nil {
		this.Ctx.Output.Status = 503
		result[constants.ERROR] = err.Error()
		response, _ := json.Marshal(result)
		this.Ctx.Output.Body(response)
		this.StopRun()
	}

}


