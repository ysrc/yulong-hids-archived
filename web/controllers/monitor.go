package controllers

import (
	"strconv"
	"yulong-hids/web/models"
)

// MonitorController old contorller
type MonitorController struct {
	BaseController
}

// GetTwenty not supported
func (c *MonitorController) GetTwenty() {
	cli := models.NewMonitor()
	ipStr := c.Ctx.Input.Param(":ip")
	startStr := c.Ctx.Input.Param(":start")
	typeStr := c.Ctx.Input.Param(":type")
	start, _ := strconv.Atoi(startStr)
	json := cli.Query(start, 20, ipStr, typeStr)
	c.Data["json"] = json
	c.ServeJSON()
	return
}

// GetAllType not suppored
func (c *MonitorController) GetAllType() {
	cli := models.NewMonitor()
	ipStr := c.Ctx.Input.Param(":ip")
	json := cli.GetAllType(ipStr)
	c.Data["json"] = json
	c.ServeJSON()
	return
}
