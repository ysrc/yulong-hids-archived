package models

import (
	"time"
	"yulong-hids/web/models/wmongo"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

// Statistics model definiton.
type Statistics struct {
	ID         bson.ObjectId `bson:"_id"      json:"_id,omitempty"`
	Uptime     time.Time     `bson:"uptime" json:"uptime,omitempty"`
	Type       string        `bson:"type" json:"type"`
	Info       string        `bson:"info" json:"info"`
	Count      int           `bson:"count" json:"count"`
	ServerList []string      `bson:"server_list" json:"server_list"`
	baseModel
}

func NewStatistics() Statistics {
	mdl := Statistics{}
	mdl.collectionName = "statistics"
	return mdl
}

func (c *Statistics) AllValue(fieldname string, limitnum int) interface{} {
	mConn := wmongo.Conn()
	defer mConn.Close()

	var list []interface{}
	collections := mConn.DB("").C(c.collectionName)
	err := collections.Find(bson.M{}).Distinct(fieldname, &list)

	if err != nil {
		beego.Error("Collections Find Distinct Error", err)
	}

	length := len(list)
	if limitnum > length {
		limitnum = length
	}

	return list[0:limitnum]
}

// Query db.getCollection('statistics').GetSortedTop
func (c *Statistics) Query(match bson.M, start int, limit int) []bson.M {
	return c.GetSortedTop(match, start, limit)
}
