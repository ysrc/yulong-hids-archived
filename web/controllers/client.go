package controllers

import (
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"gopkg.in/mgo.v2/bson"
)

// ClientController /client
type ClientController struct {
	BaseController
}

// Get http method
func (c *ClientController) Get() {
	cli := models.NewClient()
	var json interface{}

	ip := c.GetString("ip")     // client ip for monitor data
	timeout, _ := c.GetInt("t") // timeout for last seconds
	filter := c.GetString("q")  // filter for client list, find in mongodb

	// when open monitor modal dialog
	if ip != "" {
		c.Data["json"] = getMonitorData(ip, timeout)
		c.ServeJSON()
		return
	}

	paginator := c.InitPaginator()
	start, limit := paginator.ToParameter()

	var query bson.M
	if filter == "linux" {
		query = bson.M{"system": bson.M{"$regex": "^((?!windows)[\\s\\S])*$", "$options": "$i"}}
	} else if flag, exist := settings.ClientHealthTag[filter]; exist {
		query = bson.M{"health": flag}
	} else {
		query = utils.AllKeyRegexQuery(filter, cli)
	}

	json = cli.GetSortedTop(query, start, limit, "health", "ip")
	c.Data["json"] = json
	c.ServeJSON()
	return
}

func getMonitorData(ip string, t int) interface{} {
	esclient := utils.NewSession()
	var datalist interface{}
	if t != 0 {
		datalist = esclient.LastSecMonitorData(ip, t)
	} else {
		bquery := []byte(`{"from":0,"query":{"bool":{"must":[{"term":{"ip":"` + ip + `"}}]}},"size":10,"sort":[{"time":{"order":"desc"}}]}`)
		res := esclient.SearchInMonitor(bquery)
		hits := res["hits"].(map[string]interface{})
		datalist = hits["hits"]
	}
	return datalist
}
