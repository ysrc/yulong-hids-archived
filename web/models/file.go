package models

import (
	"time"
	"yulong-hids/web/models/wmongo"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type File struct {
	ID        bson.ObjectId `bson:"_id,omitempty"      json:"_id,omitempty"`
	Platform  string        `bson:"platform"       json:"platform"`
	System    string        `bson:"system"   json:"system"`
	Type      string        `bson:"type"   json:"type"`
	Hash      string        `bson:"hash" json:"hash"`
	Uptime    time.Time     `bson:"uptime"   json:"uptime"`
	baseModel `bson:",inline"`
}

func NewFile() File {
	mdl := File{}
	mdl.collectionName = "file"
	return mdl
}

func (c *File) Update() bool {
	mConn := wmongo.Conn()
	defer mConn.Close()

	collections := mConn.DB("").C("file")

	selector := bson.M{"system": c.System, "platform": c.Platform, "type": c.Type}
	c.Uptime = time.Now()
	err := collections.Update(selector, &c)

	if err == mgo.ErrNotFound {
		err = collections.Insert(c)
	} else if err != nil {
		beego.Error("File Insert Error:", err)
		return false
	}

	return true
}

func (c *File) FindOne(selector bson.M) File {
	mConn := wmongo.Conn()
	defer mConn.Close()

	var res File
	collections := mConn.DB("").C("file")

	if err := collections.Find(selector).One(&res); err != nil {
		beego.Debug("File find Error:", err)
	}

	return res
}
