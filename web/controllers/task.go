package controllers

import (
	"encoding/json"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// TaskController /task
type TaskController struct {
	BaseController
}

// Get method
func (c *TaskController) Get() {

	taskid := c.GetString("tid")
	var res interface{}

	paginator := c.InitPaginator()
	start, limit := paginator.ToParameter()

	if taskid == "" {
		cli := models.NewTask()
		res = cli.GetSortedTop(bson.M{}, start, limit, "-time")
	} else {
		cli := models.NewTaskResult()
		res = cli.GetSortedTop(bson.M{"task_id": bson.ObjectIdHex(taskid)}, start, limit, "-time")
	}

	c.Data["json"] = res
	c.ServeJSON()
	return
}

// Post method
func (c *TaskController) Post() {
	var j = models.NewTask()
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &j); err != nil {
		beego.Error("JSON Unmarshal error:", err)
		c.Data["json"] = models.NewErrorInfo(settings.AddTaskFailure)
		c.ServeJSON()
		return
	}
	if res := j.Save(); res {
		json := bson.M{
			"status": 1,
			"msg":    "添加任务成功",
			"Data":   j,
		}
		c.Data["json"] = json
	} else {
		c.Data["json"] = models.NewErrorInfo(settings.AddTaskFailure)
	}
	c.ServeJSON()
	return
}
