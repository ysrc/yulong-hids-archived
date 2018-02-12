package models

import (
	"time"
	"yulong-hids/web/models/wmongo"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

// Monitor model definiton.
type Monitor struct {
	ID   bson.ObjectId     `bson:"_id"      json:"_id,omitempty"`
	IP   string            `bson:"ip"       json:"ip"`
	Type string            `bson:"type" json:"type"`
	Data map[string]string `bson:"data" json:"data"`
	Time time.Time         `bson:"time"   json:"time,omitempty"`
	baseModel
}

func NewMonitor() Monitor {
	mdl := Monitor{}
	mdl.collectionName = "monitor"
	return mdl
}

func (c *Monitor) Query(start int, limit int, ip string, typeStr string) []Monitor {
	mConn := wmongo.Conn()
	defer mConn.Close()

	var monitorList []Monitor

	collections := mConn.DB("").C("monitor")

	querydic := bson.M{"ip": ip, "type": typeStr}

	if err := collections.Find(querydic).Sort("-time").Limit(limit).Skip(start).All(&monitorList); err != nil {
		beego.Error("Monitor Query(Find Sort Limit Skip All) Error", err)
	}

	return monitorList
}

func (c *Monitor) GetAllType(ip string) []string {
	mConn := wmongo.Conn()
	defer mConn.Close()

	collections := mConn.DB("").C("monitor")

	var result []string

	if err := collections.Find(bson.M{"ip": ip}).Distinct("type", &result); err != nil {
		beego.Error("Monitor GetAllType(Find or Distinct) Error", err)
	}

	return result
}
