package controllers

import (
	"encoding/json"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// NoticeController /notice
type NoticeController struct {
	BaseController
}

// Get method
func (c *NoticeController) Get() {

	cli := models.NewNotice()
	var res interface{}

	paginator := c.InitPaginator()
	start, limit := paginator.ToParameter()

	status := c.GetString("status")
	filter := c.GetString("q")

	if status == "learn" {
		res = cli.InfoRanking()
	} else {
		var query = make(map[string]interface{})

		if status == "dealed" {
			query["status"] = 1
		} else if status == "ignore" {
			query["status"] = 2
		} else {
			query["status"] = 0
		}

		queryor := utils.AllKeyRegexQuery(filter, cli)
		query = utils.MapUpdate(query, queryor)

		res = cli.GetSortedTop(query, start, limit, "status", "level", "-time")
	}
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// ChangeStatus POST
func (c *NoticeController) ChangeStatus() {

	cli := models.NewNotice()
	form := models.StatusForm{}

	var res bool

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &form); err != nil {
		beego.Error("ParseStatusForm(Json.Unmarshal):", err)
		c.Data["json"] = bson.M{"status": false}
		c.ServeJSON()
		return
	}
	beego.Debug("ParseStatusForm:", &form)

	if (form.Id == "all_learn") && (form.Status == 0) {
		// 关闭观察模式时，修改所有观察模式下的告警为未处理
		res = cli.LearnEnding()
	} else if form.Id == "learn" {
		// 学习模式的特殊情况
		notice := cli.FindOne(bson.M{"type": form.Type, "info": form.Info})
		beego.Debug("Notice in learn mode:", notice)
		nID := notice["_id"].(bson.ObjectId).Hex()
		res = cli.ChangeStatusbyId(nID, 2)
	} else {
		// 告警页面的逻辑
		if form.Id == "all" {

			tfaSwitch, _ := beego.AppConfig.Bool("TwoFactorAuth")
			if tfaSwitch {
				serverside := utils.GetPassword(beego.AppConfig.String("TwoFactorAuthKey"))
				clientside, err := c.GetUint32("pass")
				if err != nil || serverside != clientside {
					c.Data["json"] = bson.M{"status": false, "msg": "验证密码为空或者验证密码不正确，请重新输入双因子验证密码"}
					c.ServeJSON()
					return
				}
			}

			err := cli.UpdateAll(bson.M{"status": 0}, bson.M{"status": form.Status})
			if err != nil {
				beego.Error("Model UpdateAll", err)
				res = false
			} else {
				res = true
			}
		} else {
			res = cli.ChangeStatusbyId(form.Id, form.Status)
		}
	}

	c.Data["json"] = bson.M{"status": res}
	c.ServeJSON()
	return
}

// Delete DELETE method
func (c *NoticeController) Delete() {

	info := c.GetString("info")
	ntype := c.GetString("type")
	mgo := models.NewNotice()

	query := bson.M{"info": info, "type": ntype}
	query = utils.MapUpdate(query, settings.LearnNoticeQ)
	err := mgo.Remove(query)

	if err != nil {
		beego.Error("Notice remove: ", err)
		c.Data["json"] = bson.M{"status": 0}
	}

	c.Data["json"] = bson.M{"status": 1}
	c.ServeJSON()
	return
}
