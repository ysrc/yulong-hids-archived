package controllers

import (
	"path"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// DloadController /download
type DloadController struct {
	beego.Controller
}

// Get http method GET, uri: ?type=agent&system=linux&platform=32&action=check&access_token=AAA&hash=BBB
func (c *DloadController) Get() {

	reqtype := c.GetString("type")      // type : agent, data, daemon
	system := c.GetString("system")     // system : windows, linux
	platform := c.GetString("platform") // platform : 32, 64
	action := c.GetString("action")     // action: download, check

	// TODO: access_token := c.GetString("access_token")

	// query param while list
	if !utils.StringInSlice(system, settings.SystemArray) || !utils.StringInSlice(platform, settings.PlatformArray) || !utils.StringInSlice(reqtype, settings.TypeArray) {
		c.Data["json"] = 0
		c.ServeJSON()
		return
	}

	// check file is upgrade
	if action == "check" {
		hash := c.GetString("hash")
		filemgo := models.NewFile()
		res := filemgo.FindOne(bson.M{"system": system, "platform": platform, "type": reqtype})
		if hash == res.Hash {
			c.Data["json"] = 0
			c.ServeJSON()
		} else {
			c.Data["json"] = 1
			c.ServeJSON()
		}
		return
	}

	// download the file
	if action == "download" {
		c.Ctx.Output.Header("X-Content-Type-Options", "nosniff")
		filename := fmt.Sprintf("%s-%s-%s", system, platform, reqtype)
		filepath := path.Join(settings.FilePath, filename)
		c.Ctx.Output.Download(filepath)
		return
	}

	c.Data["json"] = 0
	c.ServeJSON()
	return

}

// Head HTTP method HEAD, bitsadmin will Head before Get
func (c *DloadController) Head() {
	c.Ctx.Output.Header("Content-Type", "text/plain; charset=utf-8")
	c.ServeJSON()
	return
}
