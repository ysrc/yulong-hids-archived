package controllers

import (
	"encoding/json"
	"yulong-hids/web/models"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

// RuleController /rules
type RuleController struct {
	BaseController
}

// Get method
func (c *RuleController) Get() {

	var res = []bson.M{}
	ruleModel := models.NewRule()

	action := c.GetString("action")
	if action == "download" {
		res = ruleModel.GetAll()
		for _, rule := range res {
			delete(rule, "_id")
		}
		c.Data["json"] = res
		c.Ctx.Output.Header("Content-Disposition", "attachment; filename=rules.json")
		c.Ctx.Output.JSON(res, true, false)

	} else {
		res = ruleModel.GetAll()
		c.Data["json"] = res
		c.ServeJSON()
	}
	return
}

// Post method
func (c *RuleController) Post() {

	var res interface{}
	ruleModel := models.NewRule()
	action := c.GetString("action")

	if action == "enable" {
		// enable or unalbe rule
		var postjson bson.M
		json.Unmarshal(c.Ctx.Input.RequestBody, &postjson)
		idstr := postjson["id"].(string)
		status := postjson["enable"].(bool)
		beego.Debug(idstr, status)
		err := ruleModel.UpdateByID(bson.ObjectIdHex(idstr), bson.M{"enabled": status})
		if err != nil {
			beego.Error("Rule change enable(model.UpdateByID):", err)
			res = bson.M{"err": err}
		} else {
			res = postjson
		}
	} else if action == "add" {
		// new rule
		var rulelist []interface{}

		json.Unmarshal(c.Ctx.Input.RequestBody, &rulelist)

		if err := ruleModel.InsertMany(rulelist); err != nil {
			beego.Error("Rule insert(model.InsertMany) error:", err)
			res = bson.M{"err": err}
		} else {
			res = rulelist
		}

	} else if action == "del" {
		// delete a rule, delete and new is edit
		var postjson bson.M
		json.Unmarshal(c.Ctx.Input.RequestBody, &postjson)
		if err := ruleModel.Remove(bson.M{"_id": bson.ObjectIdHex(postjson["id"].(string))}); err != nil {
			res = bson.M{"err": err}
		} else {
			res = bson.M{"status": true}
		}
	}

	c.Data["json"] = res
	c.ServeJSON()
	return
}
