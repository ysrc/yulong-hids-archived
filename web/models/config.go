package models

import (
	"fmt"
	"strings"
	"yulong-hids/web/models/wmongo"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type Config struct {
	Id   bson.ObjectId `bson:"_id"      json:"_id,omitempty"`
	Type string        `bson:"type"     json:"type"`
	Dic  interface{}   `bson:"dic"      json:"dic"`
	baseModel
}

func NewConfig() Config {
	mdl := Config{}
	mdl.collectionName = "config"
	return mdl
}

func (c *Config) GetAll() []bson.M {
	reslist := c.GetPieces(nil, 0, 0)
	// config of web will not be showed in web
	for index, res := range reslist {
		if res["type"] == "web" {
			reslist = append(reslist[:index], reslist[index+1:]...)
			break
		}
	}
	return reslist
}

func (c *Config) FindOne(selector bson.M) Config {
	mConn := wmongo.Conn()
	defer mConn.Close()

	var res Config
	collections := mConn.DB("").C(c.collectionName)

	if err := collections.Find(selector).One(&res); err != nil {
		beego.Warn("Config find warning:", err)
	}

	return res
}

func (c *Config) EditByID(id string, key string, value string) bool {

	if strings.Trim(value, " ") == "" {
		return false
	}

	mConn := wmongo.Conn()
	defer mConn.Close()
	mid := bson.ObjectIdHex(id)

	vresult := utils.KeyType(key, value)

	data := bson.M{"$set": bson.M{fmt.Sprintf("dic.%s", key): vresult}}
	err := mConn.DB("").C(c.collectionName).UpdateId(mid, data)
	if err == nil {
		return true
	}
	beego.Error("Config ChangeStatusbyId(model.UpdateId) Error", err)

	return false
}

func (c *Config) DelOne(id string, key string, value string) bool {

	if strings.Trim(value, " ") == "" {
		return false
	}

	mConn := wmongo.Conn()
	defer mConn.Close()

	mid := bson.ObjectIdHex(id)
	data := bson.M{"$pull": bson.M{fmt.Sprintf("dic.%s", key): value}}
	err := mConn.DB("").C(c.collectionName).UpdateId(mid, data)
	if err == nil {
		return true
	}
	beego.Error("Config ChangeStatusbyId(model.UpdateId) Error", err)

	return false
}

func (c *Config) AddOne(id string, key string, value string) bool {

	if strings.Trim(value, " ") == "" {
		return false
	}

	mConn := wmongo.Conn()
	defer mConn.Close()
	mid := bson.ObjectIdHex(id)
	data := bson.M{"$addToSet": bson.M{fmt.Sprintf("dic.%s", key): value}}
	err := mConn.DB("").C(c.collectionName).UpdateId(mid, data)
	if err == nil {
		return true
	}
	beego.Error("Config ChangeStatusbyId(model.UpdateId) Error", err)

	return false
}

// PublicKey get {"type":"server","dic":{"publickey":"return"}}
func (c *Config) PublicKey() string {
	server := c.FindOne(bson.M{"type": "server"})
	dic := server.Dic.(bson.M)
	return dic["publickey"].(string)
}
