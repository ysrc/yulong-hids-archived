package controllers

import (
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
)

// LoginController /login
type LoginController struct {
	beego.Controller
}

// Get the first page in this project
func (c *LoginController) Get() {
	c.Ctx.Output.Header("is-login-page", "true")
	c.Data["Style"] = "login-style"
	c.TplName = "login.tpl"
}

// Post HTTP method POST
func (c *LoginController) Post() {
	json := map[string]bool{"status": false}
	username := c.GetString("username")
	passwd := c.GetString("password")
	admin := beego.AppConfig.String("username")
	passhex := beego.AppConfig.String("passwordhex")
	if username == admin && utils.Md5String(passwd) == passhex {
		c.SetSession("user", admin)
		beego.Warn("User Login :", username)
		json["status"] = true
	}

	c.Data["json"] = json
	c.ServeJSON()
	return
}
