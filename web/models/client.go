package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	Id       bson.ObjectId `bson:"_id"      json:"_id,omitempty"`
	Ip       string        `bson:"ip"       json:"ip"`
	System   string        `bson:"system"   json:"system"`
	Hostname string        `bson:"hostname" json:"hostname"`
	Type     string        `bson:"type"     json:"type"`
	Uptime   time.Time     `bson:"uptime"   json:"uptime,omitempty"`
	baseModel
}

func NewClient() Client {
	mdl := Client{}
	mdl.collectionName = "client"
	return mdl
}
