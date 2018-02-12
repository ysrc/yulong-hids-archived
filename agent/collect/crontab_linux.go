// +build linux

package collect

import (
	"io/ioutil"
	"strings"
)

// GetCrontab 获取计划任务
func GetCrontab() (resultData []map[string]string) {
	//系统计划任务
	dat, err := ioutil.ReadFile("/etc/crontab")
	if err != nil {
		return resultData
	}
	cronList := strings.Split(string(dat), "\n")
	for _, info := range cronList {
		if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 6 {
			continue
		}
		s := strings.SplitN(info, " ", 7)
		rule := strings.Split(info, " "+s[5])[0]
		m := map[string]string{"command": s[6], "user": s[5], "rule": rule}
		resultData = append(resultData, m)
	}

	//用户计划任务
	dir, err := ioutil.ReadDir("/var/spool/cron/")
	if err != nil {
		return resultData
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		dat, err = ioutil.ReadFile("/var/spool/cron/" + f.Name())
		if err != nil {
			continue
		}
		cronList = strings.Split(string(dat), "\n")
		for _, info := range cronList {
			if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 5 {
				continue
			}
			s := strings.SplitN(info, " ", 6)
			rule := strings.Split(info, " "+s[5])[0]
			m := map[string]string{"command": s[5], "user": f.Name(), "rule": rule}
			resultData = append(resultData, m)
		}
	}
	return resultData
}
