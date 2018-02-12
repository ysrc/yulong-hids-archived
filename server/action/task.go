package action

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"time"
	"yulong-hids/server/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type queue struct {
	ID      bson.ObjectId `bson:"_id"`
	TaskID  bson.ObjectId `bson:"task_id"`
	IP      string        `bson:"ip"`
	Type    string        `bson:"type"`
	Command string        `bson:"command"`
	Time    time.Time     `bson:"time"`
}
type taskResult struct {
	TaskID bson.ObjectId `bson:"task_id"`
	IP     string        `bson:"ip"`
	Status string        `bson:"status" json:"status"`
	Data   string        `bson:"data" json:"data"`
	Time   time.Time     `bson:"time"`
}

var threadpool chan bool

// TaskThread 开启任务线程
func TaskThread() {
	log.Println("Start Task Thread")
	threadpool = make(chan bool, 100)
	for {
		res := queue{}
		change := mgo.Change{
			Remove: true,
		}
		models.DB.C("queue").Find(bson.M{}).Limit(1).Apply(change, &res)
		if res.IP == "" {
			time.Sleep(time.Second * 10)
			continue
		}
		threadpool <- true
		go sendTask(res, threadpool)
	}
}

func saveError(task queue, errMsg string) {
	log.Println(errMsg)
	res := taskResult{task.ID, task.IP, "false", errMsg, time.Now()}
	c := models.DB.C("task_result")
	err := c.Insert(&res)
	if err != nil {
		log.Println(err.Error())
	}
}

func sendTask(task queue, threadpool chan bool) {
	defer func() {
		<-threadpool
	}()
	sendData := map[string]string{"type": task.Type, "command": task.Command}
	if data, err := json.Marshal(sendData); err == nil {
		conn, err := net.DialTimeout("tcp", task.IP+":65512", time.Second*3)
		log.Println("sendtask:", task.IP, sendData)
		if err != nil {
			saveError(task, err.Error())
			return
		}
		defer conn.Close()
		encryptData, err := rsaEncrypt(data)
		if err != nil {
			saveError(task, err.Error())
			return
		}
		conn.Write([]byte(base64.RawStdEncoding.EncodeToString(encryptData) + "\n"))
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil || len(msg) == 0 {
			saveError(task, err.Error())
			return
		}
		log.Println(conn.RemoteAddr().String(), msg)
		res := taskResult{}
		err = json.Unmarshal([]byte(msg), &res)
		if err != nil {
			saveError(task, err.Error())
			return
		}
		res.TaskID = task.TaskID
		res.Time = time.Now()
		res.IP = task.IP
		c := models.DB.C("task_result")
		err = c.Insert(&res)
		if err != nil {
			saveError(task, err.Error())
			return
		}
	}
}
