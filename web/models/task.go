package models

import (
	"regexp"
	"strings"
	"time"
	"yulong-hids/web/models/wmongo"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	ID        bson.ObjectId `bson:"_id,omitempty"      json:"_id,omitempty"`
	Name      string        `bson:"name"       json:"name"`
	Time      time.Time     `bson:"time"   json:"time"`
	HostList  []string      `bson:"host_list"   json:"host_list"`
	Type      string        `bson:"type"   json:"type"`
	Command   string        `bson:"command"   json:"command"`
	baseModel `bson:",inline"`
}

func NewTask() Task {
	mdl := Task{}
	mdl.collectionName = "task"
	return mdl
}

func addQueue(id bson.ObjectId, c *Task, ip string) {
	queue := NewQueue()
	queue.TaskID = id
	queue.Type = c.Type
	queue.Command = c.Command
	queue.IP = ip
	if res := queue.Save(); !res {
		beego.Error("Queue insert Error, task_id:", id)
	}
}

func (c *Task) Save() bool {
	mConn := wmongo.Conn()
	defer mConn.Close()
	id := bson.NewObjectId()
	collections := mConn.DB("").C(c.collectionName)
	clientdb := NewClient()

	c.Time = time.Now()
	c.ID = id

	if err := collections.Insert(&c); err != nil {
		beego.Error("Task Insert Error", err)
		return false
	}

	// enable "all" tag
	if utils.StringInSlice("all", c.HostList) {
		tmp := clientdb.Distinct(nil, "ip")
		c.HostList = []string{}
		for _, ip := range tmp {
			c.HostList = append(c.HostList, ip.(string))
		}
	}

	// enable "windows" and "linux" tags
	if utils.StringInSlice("windows", c.HostList) {
		c.HostList = utils.DeleteElementInSlient(c.HostList, "windows")
		tmp := clientdb.Distinct(bson.M{
			"system": bson.M{"$regex": regexp.QuoteMeta("windows"), "$options": "$i"},
		}, "ip")
		for _, ip := range tmp {
			c.HostList = append(c.HostList, ip.(string))
		}
	}

	if utils.StringInSlice("linux", c.HostList) {
		c.HostList = utils.DeleteElementInSlient(c.HostList, "linux")
		tmp := clientdb.Distinct(
			bson.M{"system": bson.M{"$regex": "^((?!windows)[\\s\\S])*$", "$options": "$i"}},
			"ip",
		)
		for _, ip := range tmp {
			c.HostList = append(c.HostList, ip.(string))
		}
	}

	addedLst := []string{} // ip addresses which has added to the queue

	for _, ip := range c.HostList {
		if utils.StringInSlice(ip, addedLst) {
			break
		} else if strings.Contains(ip, "-") {
			oriLst := strings.Split(ip, "-")
			if len(oriLst) == 2 {
				subIPLst := utils.BetweenIP(oriLst[0], oriLst[1])
				for _, subIP := range subIPLst {
					addQueue(id, c, subIP)
					addedLst = append(addedLst, ip)
				}
			}
		} else {
			addQueue(id, c, ip)
			addedLst = append(addedLst, ip)
		}
	}
	return true
}

func (c *Task) ChangeStatusbyId(id string, status int) bool {
	mConn := wmongo.Conn()
	defer mConn.Close()

	mid := bson.ObjectIdHex(id)
	data := bson.M{"$set": bson.M{"status": status}}

	err := mConn.DB("").C(c.collectionName).UpdateId(mid, data)

	if err == nil {
		return true

	}

	beego.Error("Task ChangeStatusbyId(UpdateId) Error", err)

	return false
}
