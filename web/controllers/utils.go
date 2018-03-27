package controllers

import (
	"encoding/json"
	"strings"
	"time"
	"yulong-hids/web/models"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

// GetServerList Get all server url string
func GetServerList() []string {
	serverMdl := models.NewServer()
	var serverlist []string
	reslist := serverMdl.GetAll()

	for _, res := range reslist {
		serverlist = append(serverlist, res["netloc"].(string))
	}
	return serverlist
}

// GetAliveServerList return alive server url list
func GetAliveServerList() []string {
	serverMdl := models.NewServer()
	var serverlist []string

	oneMinuteAgo := time.Now().Add(time.Minute * time.Duration(-1))
	reslist := serverMdl.FindAll(bson.M{"uptime": bson.M{"$gte": oneMinuteAgo}})

	for _, res := range reslist {
		serverlist = append(serverlist, res["netloc"].(string))
	}
	return serverlist
}

// WatchModeExempt emmm...
func WatchModeExempt(c *BaseController) bool {
	// i have no choice, it's not a good way for this
	// but my bro want to make it easy...

	// checkit step by step is better then "&&" condition
	beego.Debug("Url: ", c.Ctx.Input.URL())
	if !strings.HasSuffix(c.Ctx.Input.URL(), "config") {
		return false
	}

	mgocli := models.NewConfig()
	serverConf := mgocli.FindOne(bson.M{"type": "server"})
	isLearn := serverConf.Dic.(bson.M)["learn"].(bool)
	if !isLearn {
		return false
	}
	whiteID := mgocli.FindOne(bson.M{"type": "whitelist"}).Id
	editf := models.EditCfgForm{}
	if json.Unmarshal(c.Ctx.Input.RequestBody, &editf); editf.Id != whiteID.Hex() {
		return false
	}
	return true
}

// IsLearn is watching mode
func IsLearn() bool {
	mgocli := models.NewConfig()
	serverConf := mgocli.FindOne(bson.M{"type": "server"})
	is := serverConf.Dic.(bson.M)["learn"].(bool)
	return is
}
