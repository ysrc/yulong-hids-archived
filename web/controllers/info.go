package controllers

import (
	"yulong-hids/web/models"

	"gopkg.in/mgo.v2/bson"
)

// InfoController /info
type InfoController struct {
	BaseController
}

// Get not route to this controller
func (c *InfoController) Get() {
	cli := models.NewInfo()
	json := cli.GetAll()
	c.Data["json"] = json
	c.ServeJSON()
	return
}

// GetInfoByIp HTTP method GET
func (c *InfoController) GetInfoByIp() {
	ipStr := c.Ctx.Input.Param(":ip")
	info := models.NewInfo()
	client := models.NewClient()

	infoData := info.GetInfoByIp(ipStr)
	clientData := client.FindOne(bson.M{"ip": ipStr})

	clientData["infodata"] = infoData
	c.Data["json"] = clientData
	c.ServeJSON()
	return
}
