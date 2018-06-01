package safecheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"yulong-hids/server/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type stats struct {
	ServerList []string `bson:"server_list"`
	Count      int      `bson:"count"`
}

// ScanChan 待检测队列
var ScanChan = make(chan models.DataInfo, 4096)
var cache []string

// Check 检测引擎结构
type Check struct {
	Info        models.DataInfo   // 待检测数据集
	V           map[string]string // 当前检测数据内容
	Value       string            // 触发规则的信息
	Description string            // 规则简介信息
	Source      string            // 警报来源
	Level       int               // 警报等级
	CStatistics *mgo.Collection   // 统计表
	CNoice      *mgo.Collection   // 警报表
}

// BlackFilter 黑名单检测
func (c *Check) BlackFilter() {
	var keyword string
	var blackList []string
	regex := true
	switch c.Info.Type {
	case "process":
		blackList = models.Config.BlackList.Process
		keyword = c.V["name"]
	case "connection":
		blackList = models.Config.BlackList.IP
		keyword = strings.Split(c.V["remote"], ":")[0]
		regex = false
	case "loginlog":
		blackList = models.Config.BlackList.IP
		keyword = c.V["remote"]
		regex = false
	case "file":
		blackList = models.Config.BlackList.File
		keyword = c.V["hash"]
		regex = false
	case "crontab":
		blackList = models.Config.BlackList.File
		keyword = c.V["command"]
	default:
		blackList = models.Config.BlackList.Other
		keyword = c.V["name"]
	}
	if len(blackList) >= 1 && inArray(blackList, strings.ToLower(keyword), regex) {
		c.Source = "blacklist"
		c.Level = 0
		c.Description = "存在于黑名单列表中"
		c.Value = keyword
		c.warning()
	}
}

// WhiteFilter 白名单筛选
func (c *Check) WhiteFilter() bool {
	var keyword string
	var whiteList []string
	regex := true
	switch c.Info.Type {
	case "process":
		whiteList = models.Config.WhiteList.Process
		keyword = c.V["name"]
	case "connection":
		whiteList = models.Config.WhiteList.IP
		keyword = strings.Split(c.V["remote"], ":")[0]
		regex = false
	case "loginlog":
		whiteList = models.Config.WhiteList.IP
		keyword = c.V["remote"]
		regex = false
	case "file":
		whiteList = models.Config.WhiteList.File
		keyword = c.V["hash"]
		regex = false
	case "crontab":
		whiteList = models.Config.WhiteList.File
		keyword = c.V["command"]
	default:
		whiteList = models.Config.WhiteList.Other
		keyword = c.V["name"]
	}
	if len(whiteList) >= 1 && inArray(whiteList, strings.ToLower(keyword), regex) {
		return true
	}
	return false
}

// Rules 对预定规则解析匹配
func (c *Check) Rules() {
	for _, r := range models.RuleDB {
		var vulInfo []string
		if (c.Info.System != r.System && r.System != "all") || c.Info.Type != r.Source {
			continue
		}
		i := len(r.Rules)
		// log.Println(r.Rules)
		for k, rule := range r.Rules {
			switch rule.Type {
			case "regex":
				reg := regexp.MustCompile(rule.Data)
				if reg.MatchString(strings.ToLower(c.V[k])) {
					i--
					vulInfo = append(vulInfo, c.V[k])
				}
			case "non-regex":
				reg := regexp.MustCompile(rule.Data)
				if c.V[k] != "" && !reg.MatchString(strings.ToLower(c.V[k])) {
					i--
					vulInfo = append(vulInfo, c.V[k])
				}
			case "string":
				if strings.ToLower(c.V[k]) == strings.ToLower(rule.Data) {
					i--
					vulInfo = append(vulInfo, c.V[k])
				}
			case "count":
				if models.Config.Learn {
					i--
					vulInfo = append(vulInfo, c.V[k])
					continue
				}
				var statsinfo stats
				var keyword string
				if c.Info.Type == "connection" {
					keyword = strings.Split(c.V[k], ":")[0]
				} else {
					keyword = c.V[k]
				}
				err := c.CStatistics.Find(bson.M{"type": r.Source, "info": keyword}).One(&statsinfo)
				if err != nil {
					log.Println(err.Error(), r.Source, keyword)
					break
				}
				n, err := strconv.Atoi(rule.Data)
				if err != nil {
					log.Println(err.Error())
					break
				}
				if statsinfo.Count == n {
					i--
					vulInfo = append(vulInfo, c.V[k])
				}
			}
		}
		if r.And {
			if i == 0 {
				c.Source = r.Meta.Name
				c.Level = r.Meta.Level
				c.Description = r.Meta.Description
				sort.Strings(vulInfo)
				c.Value = strings.Join(vulInfo, "|")
				c.warning()
			}
		} else if i < len(r.Rules) {
			c.Source = r.Meta.Name
			c.Level = r.Meta.Level
			c.Description = r.Meta.Description
			sort.Strings(vulInfo)
			c.Value = strings.Join(vulInfo, "|")
			c.warning()
		}
	}
}

// Intelligence 威胁情报接口检测
func (c *Check) Intelligence() {
	if !models.Config.Intelligence.Switch {
		return
	}
	if c.Info.Type == "connection" || c.Info.Type == "loginlog" {
		ip := strings.Split(c.V["remote"], ":")[0]
		if isLan(ip) {
			return
		}
		if inArray(cache, c.Info.Type+c.Info.IP+ip, false) {
			return
		}
		cache = append(cache, c.Info.Type+c.Info.IP+ip)
		url := strings.Replace(models.Config.Intelligence.IPAPI, "{$ip}", ip, 1)
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		reg := regexp.MustCompile(models.Config.Intelligence.Regex)
		if reg.Match(body) {
			c.Source = "威胁情报接口"
			c.Level = 0
			c.Description = "威胁情报接口显示此IP存在风险"
			c.Value = ip
			c.warning()
		}
	} else if c.Info.Type == "file" {
		if c.V["hash"] != "" {
			if inArray(cache, c.Info.Type+c.Info.IP+c.V["hash"], false) {
				return
			}
			cache = append(cache, c.Info.Type+c.Info.IP+c.V["hash"])
			url := strings.Replace(models.Config.Intelligence.FileAPI, "{$hash}", c.V["hash"], 1)
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			reg := regexp.MustCompile(models.Config.Intelligence.Regex)
			if reg.Match(body) {
				c.Source = "威胁情报接口"
				c.Level = 0
				c.Description = "威胁情报接口显示此文件存在风险"
				c.Value = c.V["hash"]
				c.warning()
			}
		}
	}
}

func (c *Check) warning() {
	// 观察模式 只记录统计 不显示
	if models.Config.Learn {
		c.CNoice.Upsert(bson.M{"type": c.Info.Type, "ip": c.Info.IP, "source": c.Source, "level": c.Level,
			"info": c.Value, "description": c.Description, "status": 3, "time": c.Info.Uptime}, bson.M{"$inc": bson.M{"status": 1}})
	} else {
		// 如果忽略过就不写入
		n, _ := c.CNoice.Find(bson.M{"type": c.Info.Type, "ip": c.Info.IP, "info": c.Value, "status": 2}).Count()
		if n >= 1 {
			return
		}
		raw, err := json.Marshal(c.V)
		if err != nil {
			log.Println(err.Error())
		}
		err = c.CNoice.Insert(bson.M{"type": c.Info.Type, "ip": c.Info.IP, "source": c.Source, "level": c.Level,
			"info": c.Value, "description": c.Description, "status": 0, "raw": string(raw), "time": c.Info.Uptime})
		if err == nil {
			msg := fmt.Sprintf("IP:%s,Type:%s,Info:%s %s", c.Info.IP, c.Info.Type, c.Value, c.Description)
			sendNotice(c.Level, msg)
		}
	}
}

// Run 开始检查
func (c *Check) Run() {
	for _, c.V = range c.Info.Data {
		c.BlackFilter()
		if c.WhiteFilter() {
			continue
		}
		c.Rules()
		c.Intelligence()
	}
}

// ScanMonitorThread 安全检测线程
func ScanMonitorThread() {
	log.Println("Start Scan Thread")
	// 10个检测goroutine
	for i := 0; i < 10; i++ {
		go func() {
			c := new(Check)
			c.CStatistics = models.DB.C("statistics")
			c.CNoice = models.DB.C("notice")
			for {
				c.Info = <-ScanChan
				c.Run()
			}
		}()
	}
	ticker := time.NewTicker(time.Second * 60)
	for _ = range ticker.C {
		cache = []string{}
	}
}
