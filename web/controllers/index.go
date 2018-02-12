package controllers

import (
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
)

// MainController index url:/
type MainController struct {
	BaseController
}

// Get HTTP method GET
func (c *MainController) Get() {
	if HasInstall() < settings.InstallStep {
		c.Ctx.Redirect(302, beego.URLFor("InstallController.Get"))
		return
	}

	token := utils.RandStringBytesMaskImprSrc(64) // csrf token

	isTwoFactorAuth, _ := beego.AppConfig.Bool("TwoFactorAuth")

	c.Data["version"] = settings.Version
	c.Data["iswatchmode"] = IsLearn()
	c.Data["apihost"] = beego.AppConfig.String("apihost")
	c.Data["httpport"] = beego.AppConfig.String("HTTPPort")
	c.Data["token"] = token
	c.Data["istwofactorauth"] = isTwoFactorAuth
	c.Ctx.SetCookie("request_token", token, 1073741824, "/")
	c.TplName = "index.tpl"
}
