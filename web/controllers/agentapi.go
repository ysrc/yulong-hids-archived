package controllers

import "net/url"

import (
	"strings"
	"yulong-hids/web/models"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// AgentApiController web api for agent
type AgentApiController struct {
	beego.Controller
}

// Get agent will get publickey content and serverlist here
func (c *AgentApiController) Get() {

	currentURL := c.Ctx.Request.RequestURI

	if strings.Contains(currentURL, "publickey") {
		conf := models.NewConfig()
		c.Data["json"] = bson.M{"public": conf.PublicKey()}
	}

	if strings.Contains(currentURL, "serverlist") {
		c.Data["json"] = GetAliveServerList()
	}

	if strings.Contains(currentURL, "dbinfo") {
		esurl := beego.AppConfig.String("elastic_search::baseurl")
		mgourl := beego.AppConfig.String("mongodb::url")
		u, _ := url.Parse(mgourl)
		c.Data["json"] = bson.M{
			"elastic_search_url": esurl,
			"mongodb_url":        u.Host,
		}
	}

	c.ServeJSON()
	return

}
