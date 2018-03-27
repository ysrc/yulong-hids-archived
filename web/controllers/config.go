package controllers

import (
	"encoding/json"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// ConfigController /config
type ConfigController struct {
	BaseController
}

// Get HTTP method GET
func (c *ConfigController) Get() {
	cli := models.NewConfig()
	json := cli.GetAll()
	for _, item := range json {
		if item["type"] == "server" {
			dic := item["dic"].(bson.M)
			for _, secretkey := range settings.SecretKeyLst {
				dic[secretkey] = utils.Md5String(dic[secretkey].(string))
			}
		}
	}
	c.Data["json"] = json
	c.ServeJSON()
	return
}

// Edit HTTP method POST
func (c *ConfigController) Edit() {

	cli := models.NewConfig()
	j := models.EditCfgForm{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &j); err != nil {
		beego.Debug("Config edit error:", err)
		c.Data["json"] = models.NewErrorInfo(settings.EditCfgFailure)
		c.ServeJSON()
		return
	}

	res := cli.EditByID(j.Id, j.Key, j.Input)
	c.Data["json"] = bson.M{"status": res}
	c.ServeJSON()
	return
}

// Del HTTP method DELETE
func (c *ConfigController) Del() {

	cli := models.NewConfig()
	j := models.EditCfgForm{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &j); err != nil {
		beego.Debug("Config edit error:", err)
		c.Data["json"] = models.NewErrorInfo(settings.EditCfgFailure)
		c.ServeJSON()
		return
	}
	res := cli.DelOne(j.Id, j.Key, j.Input)
	c.Data["json"] = bson.M{"status": res}
	c.ServeJSON()
	return
}

// Add HTTP method PUT
func (c *ConfigController) Add() {

	cli := models.NewConfig()
	j := models.EditCfgForm{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &j); err != nil {
		beego.Debug("Config add error:", err)
		c.Data["json"] = models.NewErrorInfo(settings.EditCfgFailure)
		c.ServeJSON()
		return
	}

	// support blacklist and whitelist convenient tag
	if j.Id == "blacklist" || j.Id == "whitelist" {
		config := cli.FindOne(bson.M{"type": j.Id})
		j.Id = config.Id.Hex()
	}

	res := cli.AddOne(j.Id, j.Key, j.Input)
	c.Data["json"] = bson.M{"status": res}
	c.ServeJSON()
	return
}
