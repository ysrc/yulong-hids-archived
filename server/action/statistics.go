package action

import (
	"yulong-hids/server/models"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// ResultStat 对接收数据进行统计
func ResultStat(datainfo models.DataInfo) error {
	var err error
	c := models.DB.C("statistics")
	mainMapping := map[string]string{
		"process":    "name",
		"userlist":   "name",
		"listening":  "address",
		"connection": "remote",
		"loginlog":   "remote",
		"startup":    "name",
		"crontab":    "command",
		"service":    "name",
		// "processlist": "name",
	}
	if _, ok := mainMapping[datainfo.Type]; !ok {
		return nil
	}
	k := mainMapping[datainfo.Type]
	ip := datainfo.IP
	for _, v := range datainfo.Data {
		if datainfo.Type == "connection" {
			v[k] = strings.Split(v[k], ":")[0]
		}
		count, _ := c.Find(bson.M{"info": v[k], "type": datainfo.Type}).Count()
		if count >= 1 {
			err = c.Update(bson.M{"info": v[k], "type": datainfo.Type}, bson.M{
				"$set":      bson.M{"uptime": datainfo.Uptime},
				"$inc":      bson.M{"count": 1},
				"$addToSet": bson.M{"server_list": ip}})
		} else {
			serverList := []string{ip}
			err = c.Insert(bson.M{"type": datainfo.Type, "info": v[k], "count": 1,
				"server_list": serverList, "uptime": datainfo.Uptime})
		}
	}
	return err
}
