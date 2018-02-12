package models

import (
	"time"
	"yulong-hids/web/models/wmongo"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type Queue struct {
	ID        bson.ObjectId `bson:"_id,omitempty"      json:"_id,omitempty"`
	TaskID    bson.ObjectId `bson:"task_id"      json:"_task_id"`
	Time      time.Time     `bson:"time"   json:"time"`
	IP        string        `bson:"ip"   json:"ip"`
	Type      string        `bson:"type"   json:"type"`
	Command   string        `bson:"command"   json:"command"`
	baseModel `bson:",inline"`
}

func NewQueue() Queue {
	mdl := Queue{}
	mdl.collectionName = "queue"
	return mdl
}

func (c *Queue) Save() bool {
	mConn := wmongo.Conn()
	defer mConn.Close()

	collections := mConn.DB("").C(c.collectionName)
	c.Time = time.Now()
	if err := collections.Insert(&c); err != nil {
		beego.Error("Queue Insert Error", err)
		return false
	}
	return true
}
