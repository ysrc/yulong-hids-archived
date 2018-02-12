package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// LogoutController /login
type LogoutController struct {
	beego.Controller
}

// Post HTTP method POST
func (c *LogoutController) Post() {
	c.DelSession("user")
	c.Data["json"] = bson.M{"status": true}
	c.ServeJSON()
	return
}
