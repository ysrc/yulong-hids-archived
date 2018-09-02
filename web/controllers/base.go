package controllers

import (
	"strings"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// BaseController base web controller, base struct
type BaseController struct {
	beego.Controller
	IsRoot bool
}

// Prepare access Control, 2FA, csrf check and other security options
func (c *BaseController) Prepare() {

	// check hostname
	hostname := beego.AppConfig.String("ylhostname")
	allowHosts := strings.Split(hostname, ",")
	if hostname != "" && !utils.StringInSlice(c.Ctx.Input.Host(), allowHosts) {
		beego.Error("Hostname not correct.")
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "Forbidden"
		c.ServeJSON()
		return
	}

	// only https be allowed
	HTTPSOnly, _ := beego.AppConfig.Bool("OnlyHTTPS")
	if HTTPSOnly && c.Ctx.Input.Scheme() != "https" &&
		utils.FindSub(settings.HTTPURLLst, c.Ctx.Input.URL()) == "" {

		c.Ctx.Redirect(302, "https://"+c.Ctx.Input.Domain()+":"+beego.AppConfig.String("HTTPSPort"))
		return
	}

	// user login check
	if !utils.IsDevMode() {
		username := c.GetSession("user")
		if username != beego.AppConfig.String("username") {
			c.Ctx.Redirect(302, beego.URLFor("LoginController.Get"))
			return
		}
	}

	// two factor auth check
	tfaSwitch, _ := beego.AppConfig.Bool("TwoFactorAuth")
	var first uint32
	if c.Ctx.Input.Method() != "GET" &&
		tfaSwitch && !WatchModeExempt(c) &&
		utils.FindSub(settings.AuthURILst, c.Ctx.Input.URL()) != "" {
		serverside := utils.GetPassword(beego.AppConfig.String("TwoFactorAuthKey"))
		beego.Debug("GetPassword: ", serverside)
		// push
		settings.TFAPassHistorys = append(settings.TFAPassHistorys, serverside)
		// pop
		first, settings.TFAPassHistorys = settings.TFAPassHistorys[0], settings.TFAPassHistorys[1:]
		// only allow to try 6 times
		if first == serverside {
			c.Data["json"] = bson.M{"status": false, "msg": "尝试次数太多，请等待30秒后再进行验证"}
			c.ServeJSON()
			return
		}
		clientside, err := c.GetUint32("pass")
		if err != nil || serverside != clientside {
			c.Data["json"] = bson.M{"status": false, "msg": "验证密码为空或者验证密码不正确，请重新输入双因子验证密码"}
			c.ServeJSON()
			return
		}
	}

	// csrf check protection
	if c.Ctx.Input.Method() != "GET" && !utils.IsDevMode() {
		cookietoken := c.Ctx.GetCookie("request_token")
		headertoken := c.Ctx.Input.Header("RequestToken")
		if headertoken != cookietoken {
			c.Data["json"] = bson.M{"status": false, "msg": "CSRF检测错误，请刷新页面重试。"}
			c.ServeJSON()
			return
		}
	}

	// set security http headers
	c.Ctx.Output.Header("X-Frame-Options", "SAMEORIGIN")
	c.Ctx.Output.Header("X-Content-Type-Options", "nosniff")
	c.Ctx.Output.Header("X-XSS-Protection", "1; mode=block")
	c.Ctx.Output.Header("X-Robots-Tag", "none")
	c.Ctx.Output.Header("X-Download-Options", "noopen")
	c.Ctx.Output.Header("X-Permitted-Cross-Domain-Policies", "none")

	// TODO dev mode to dev-vul model, you never set it
	if beego.AppConfig.String("runmode") == "dev-vul" {
		c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
		c.Ctx.Output.Header("Content-Security-Policy", "")
		c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	}

	beego.Info("Url:", c.Ctx.Request.RequestURI)
}

// Options : chrome "preflighted" requests first send an HTTP request by the OPTIONS method to the resource on the other domain
func (c *BaseController) Options() {
	c.Data["json"] = bson.M{"status": false, "msg": "allow browers look at my CORS HEADERs"}
	c.ServeJSON()
	return
}

// InitPaginator : Paginator function
func (c *BaseController) InitPaginator() utils.Paginator {
	pagenum, _ := c.GetInt("page")
	limit, limiterr := c.GetInt("limit")

	if limiterr != nil {
		limit = settings.PageLimit
	}

	p := utils.Paginator{}
	p.Pagenum = pagenum
	p.Limit = limit

	return p
}
