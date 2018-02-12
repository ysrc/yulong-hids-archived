package models

import (
	"encoding/json"
	"time"
	"yulong-hids/web/models/wmongo"
	"yulong-hids/web/settings"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

// User model definiton.
type Notice struct {
	Id     bson.ObjectId `bson:"_id"      json:"_id,omitempty"`
	Info   string        `bson:"info"       json:"info"`
	Status int           `bson:"status"   json:"status"`
	Time   time.Time     `bson:"time" json:"time,omitempty"`
	Type   string        `bson:"type"   json:"type"`
	Ip     string        `bson:"ip"   json:"ip"`
	Source string        `bson:"source"   json:"source"`
	Level  int           `bson:"level"   json:"level"`
	baseModel
}

func NewNotice() Notice {
	mdl := Notice{}
	mdl.collectionName = "notice"
	return mdl
}

func (c *Notice) ChangeStatusbyId(id string, status int) bool {
	mConn := wmongo.Conn()
	defer mConn.Close()

	mid := bson.ObjectIdHex(id)
	data := bson.M{"$set": bson.M{"status": status, "uptime": time.Now()}}

	err := mConn.DB("").C(c.collectionName).UpdateId(mid, data)

	if err == nil {
		return true
	}

	beego.Error("Notice ChangeStatusbyId(UpdateId) Error", err)

	return false
}

func (c *Notice) CountPerByKey(match bson.M, key string) []bson.M {
	var sort, group bson.M
	jsonGroup := []byte(`{"$group": {"_id": "$` + key + `", "count": {"$sum": 1}}}`)
	jsonSort := []byte(`{"$sort": {"count": -1}}`)
	json.Unmarshal(jsonGroup, &group)
	json.Unmarshal(jsonSort, &sort)
	res := c.Aggregate(match, group, sort)
	return res
}

func (c *Notice) CountPerDay(match bson.M) []bson.M {
	var project, group bson.M
	json.Unmarshal(settings.StatisticsPipeProjectQ, &project)
	json.Unmarshal(settings.StatisticsPipeGroupQ, &group)
	res := c.Aggregate(bson.M{"$match": match}, project, group)
	return res
}

func (c *Notice) CountPerHour(match bson.M) []bson.M {
	project := bson.M{
		"$project": bson.M{
			"d":    bson.M{"$dateToString": bson.M{"format": "%m-%d:%H", "date": "$time"}},
			"time": 1,
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "$d",
			"count": bson.M{
				"$sum": 1,
			},
			"time": bson.M{
				"$first": "$time",
			},
		},
	}
	sort := bson.M{"$sort": bson.M{"_id": 1}}
	res := c.Aggregate(bson.M{"$match": match}, project, group, sort)
	return res
}

func (c *Notice) InfoRanking() []bson.M {

	match := bson.M{
		"type":   bson.M{"$in": settings.Learn2WriteList},
		"status": bson.M{"$gt": 2},
	}

	group := bson.M{
		"$group": bson.M{
			"_id":   bson.M{"info": "$info", "type": "$type"},
			"count": bson.M{"$sum": 1},
		},
	}
	sort := bson.M{"$sort": bson.M{"count": -1}}
	res := c.Aggregate(bson.M{"$match": match}, group, sort)
	return res
}

// UpdateAll update many
func (c *Notice) UpdateAll(selector bson.M, update bson.M) error {

	mConn := wmongo.Conn()
	defer mConn.Close()
	cname := c.collectionName
	collection := mConn.DB("").C(cname)

	update = bson.M{"$set": update}
	_, err := collection.UpdateAll(selector, update)

	return err

}

// LearnEnding done watch mode ending action
func (c *Notice) LearnEnding() bool {

	selector := settings.LearnNoticeQ
	update := bson.M{"status": 0}

	err := c.UpdateAll(selector, update)
	if err != nil {
		beego.Error("Notice UpdateAll:", err)
		go c.doneDuplicateErr()
	}

	return true
}

//处理因为index而产生的问题数据
func (c *Notice) doneDuplicateErr() {

	selector := settings.LearnNoticeQ
	update := bson.M{"$set": bson.M{"status": 0}}

	allColl := c.FindAll(selector)
	for _, coll := range allColl {
		err := c.UpdateByID(coll["_id"], update)
		if err != nil {
			c.Remove(selector)
		}
	}
}
