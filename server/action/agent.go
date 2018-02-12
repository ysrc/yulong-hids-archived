package action

import (
	"log"
	"time"
	"yulong-hids/server/models"

	"gopkg.in/mgo.v2/bson"
)

type client struct {
	TYPE string
	DIC  ClientConfig
}
type monitorInfo struct {
	IP   string
	Type string
	Time time.Time
}

// ClientConfig 客户端配置信息结构
type ClientConfig struct {
	Cycle       int      `bson:"cycle"` // 信息传输频率，单位：分钟
	UDP         bool     `bson:"udp"`   // 是否记录UDP请求
	LAN         bool     `bson:"lan"`   // 是否本地网络请求
	Mode        string   `bson:"mode"`  // 模式，考虑中
	Filter      filter   // 直接过滤不回传的数据
	MonitorPath []string `bson:"monitorPath"` // 监控目录列表
	Lasttime    string   // 最后一条登录日志时间
}
type filterres struct {
	Type string `bson:"type"`
	Dic  filter `bson:"dic"`
}
type filter struct {
	File    []string `bson:"file"`    // 文件hash、文件名
	IP      []string `bson:"ip"`      // IP地址
	Process []string `bson:"process"` // 进程名、参数
}

// GetAgentConfig 返回客户端的配置信息
func GetAgentConfig(ip string) ClientConfig {
	var clientRes client
	c := models.DB.C("config")
	c.Find(bson.M{"type": "client"}).One(&clientRes)
	config := clientRes.DIC

	var res filterres
	c.Find(bson.M{"type": "filter"}).One(&res)
	config.Filter = res.Dic
	lastTime, err := models.QueryLogLastTime(ip)
	if err != nil {
		log.Println(err.Error())
		config.Lasttime = "all"
	} else {
		config.Lasttime = lastTime
	}
	return config
}
