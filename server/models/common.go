package models

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// DB 数据库连接池
	DB      *mgo.Database
	mongodb *string
	es      *string
	// Config 配置信息
	Config serverConfig
	// LocalIP 本机活动IP
	LocalIP string
	err     error
	// RuleDB 存放在mongodb rule 的规则库
	RuleDB = []ruleInfo{}
)

// DataInfo 从agent接收数据的结构
type DataInfo struct {
	IP     string              `json:"ip"`
	Type   string              `json:"type"`
	System string              `json:"system"`
	Data   []map[string]string `json:"data"`
	Uptime time.Time           `json:"uptime"`
}

type configres struct {
	Type string       `bson:"type"`
	Dic  serverConfig `bson:"dic"`
}
type intelligencegres struct {
	Type string       `bson:"type"`
	Dic  intelligence `bson:"dic"`
}
type intelligence struct {
	Switch  bool   `bson:"switch"`  // 开关
	IPAPI   string `bson:"ipapi"`   // 查询IP是否存在风险的URL接口
	FileAPI string `bson:"fileapi"` // 查询文件md5是否存在风险的URL接口
	Regex   string `bson:"regex"`   // 判断为威胁的正则表达式
}
type noticeres struct {
	Type string `bson:"type"`
	Dic  notice `bson:"dic"`
}
type notice struct {
	Switch   bool   `bson:"switch"`   // 开关
	API      string `bson:"api"`      // API URL接口
	OnlyHigh bool   `bson:"onlyhigh"` // 仅通知危险等级的告警
}
type blackListres struct {
	Type string    `bson:"type"`
	Dic  blackList `bson:"dic"`
}
type blackList struct {
	File    []string `bson:"file"`    // 文件hash值
	IP      []string `bson:"ip"`      // IP地址
	Process []string `bson:"process"` // 进程名称或者完整命令
	Other   []string `bson:"other"`   // 其他name
}
type whiteListres struct {
	Type string    `bson:"type"`
	Dic  whiteList `bson:"dic"`
}
type whiteList struct {
	File    []string `bson:"file"`    // 文件hash、文件名
	IP      []string `bson:"ip"`      // IP地址
	Process []string `bson:"process"` // 进程名、参数
	Other   []string `bson:"other"`   // 其他name
}

type serverConfig struct {
	Learn        bool         `bson:"learn"`        // 是否为观察模式
	OfflineCheck bool         `bson:"offlinecheck"` // 开启离线主机检测和通知
	BlackList    blackList    // 黑名单
	WhiteList    whiteList    // 白名单
	Private      string       `bson:"privatekey"` // 加密秘钥
	Cert         string       `bson:"cert"`       // TLS加密证书
	Intelligence intelligence // 威胁情报
	Notice       notice       // 通知
}

type rule struct {
	Type string `json:"type" bson:"type"`
	Data string `json:"data" bson:"data"`
}

type ruleInfo struct {
	Meta struct {
		Name        string `json:"name" bson:"name"`               // 名称
		Author      string `json:"author" bson:"author"`           // 编写人
		Description string `json:"description" bson:"description"` // 描述
		Level       int    `json:"level" bson:"level"`             // 风险等级
	} `json:"meta" bson:"meta"` // 规则信息
	Source string          `json:"source" bson:"source"` // 选择判断来源
	System string          `json:"system" bson:"system"` // 匹配系统
	Rules  map[string]rule `json:"rules" bson:"rules"`   // 具体匹配规则
	And    bool            `json:"and" bson:"and"`       // 规则逻辑
}

func init() {
	mongodb = flag.String("db", "", "mongodb ip:port")
	es = flag.String("es", "", "elasticsearch ip:port")
	flag.Parse()
	if len(os.Args) <= 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.HasPrefix(*mongodb, "127.") || strings.HasPrefix(*mongodb, "localhost") {
		log.Println("mongodb Can not be 127.0.0.1")
		os.Exit(1)
	}
	DB, err = conn(*mongodb, "agent")
	if err != nil {
		log.Println(err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}
	LocalIP, err = getLocalIP(*mongodb)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Get Config")
	setConfig()
	setRules()
	go esCheckThread()
}
func getLocalIP(ip string) (string, error) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

// setConfig 获取配置文件
func setConfig() {
	c := DB.C("config")
	res := configres{}
	c.Find(bson.M{"type": "server"}).One(&res)

	res2 := intelligencegres{}
	c.Find(bson.M{"type": "intelligence"}).One(&res2)

	res3 := blackListres{}
	c.Find(bson.M{"type": "blacklist"}).One(&res3)

	res4 := whiteListres{}
	c.Find(bson.M{"type": "whitelist"}).One(&res4)

	res5 := noticeres{}
	c.Find(bson.M{"type": "notice"}).One(&res5)

	Config = res.Dic
	Config.Intelligence = res2.Dic
	Config.BlackList = res3.Dic
	Config.WhiteList = res4.Dic
	Config.Notice = res5.Dic
}

// setRules 获取异常规则集
func setRules() {
	c := DB.C("rules")
	c.Find(bson.M{"enabled": true}).All(&RuleDB)
}

// regServer 注册为服务，Agent才知道发给谁
func regServer() {
	c := DB.C("server")
	_, err := c.Upsert(bson.M{"netloc": LocalIP + ":33433"}, bson.M{"$set": bson.M{"uptime": time.Now()}})
	if err != nil {
		log.Println(err.Error())
	}
}

// Heartbeat 心跳线程，定时刷新配置和规则
func Heartbeat() {
	log.Println("Start heartbeat thread")
	for {
		mgoCheck()
		regServer()
		setConfig()
		setRules()
		time.Sleep(time.Second * 30)
	}
}
