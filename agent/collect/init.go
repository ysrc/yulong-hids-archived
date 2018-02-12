// Package collect 获取以下服务器关键信息
// 监听端口，服务列表，用户列表，启动项，计划任务，登录日志
package collect

import (
	"regexp"
	"time"
	"yulong-hids/agent/common"
)

var allInfo = make(map[string][]map[string]string)

var tagMap = map[string]string{
	"web": `nginx|httpd|apache|w3wp\.exe|tomcat|weblogic|jboss|jetty`,
	"db":  `mysql|mongo|sqlservr\.exe|oracle|elasticsearch|postgres|redis|cassandra|teradata|solr|HMaster|hbase|mariadb`,
}

func init() {
	go func() {
		time.Sleep(time.Second * 3600)
		common.ServerInfo = GetComInfo()
	}()
}

// GetAllInfo 获取所有收集的信息
func GetAllInfo() map[string][]map[string]string {
	allInfo["listening"] = GetListening()
	allInfo["service"] = GetServiceInfo()
	allInfo["userlist"] = GetUser()
	allInfo["startup"] = GetStartup()
	allInfo["crontab"] = GetCrontab()
	allInfo["loginlog"] = GetLoginLog()
	allInfo["processlist"] = GetProcessList()
	return allInfo
}

func discern(info *common.ComputerInfo) {
	for k, v := range tagMap {
		for _, p := range GetProcessList() {
			if p["command"] == "" {
				continue
			}
			if ok, _ := regexp.MatchString(v, p["command"]); ok {
				info.Type = k
				if k == "web" {
					info.Path, _ = getWebPath(p["command"])
					// web优先，匹配到web就退出，其他一直匹配下去
					return
				}
			}
		}
	}
}

func removeDuplicatesAndEmpty(list []string) (ret []string) {
	listLen := len(list)
	for i := 0; i < listLen; i++ {
		if (i > 0 && list[i-1] == list[i]) || len(list[i]) == 0 {
			continue
		}
		ret = append(ret, list[i])
	}
	return
}
