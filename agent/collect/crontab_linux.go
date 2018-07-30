// +build linux

package collect

import (
	"io/ioutil"
	"strings"
)

// GetCrontab 获取计划任务
/*func GetCrontab() (resultData []map[string]string) {
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
}*/

//修复crontab中空格不规则导致异常
func GetCrontab() (resultData []map[string]string) {
	//系统计划任务
	dat, err := ioutil.ReadFile("/etc/crontab")
	if err != nil {
		d <- resultData
		return resultData
	}
	cronList := strings.Split(string(dat), "\n")
	for _, info := range cronList {
		if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 6 {
			continue
		}
		//s := strings.SplitN(info, " ", 7)
		_, cmd := GetCmd(info, 7)
		tmp := strings.Split(info, " "+cmd)[0]
		_, user := GetCmd(tmp, 6)
		rule := strings.Split(tmp, " "+user)[0]
		m := map[string]string{"command": cmd, "user": user, "rule": rule}
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
			_, cmd := GetCmd(info, 6)
			rule := strings.Split(info, " "+cmd)[0]
			m := map[string]string{"command": cmd, "user": f.Name(), "rule": rule}
			resultData = append(resultData, m)
		}
	}
	return resultData
}

func GetCmd(r string, c int) (rule string, cmd string) {

	n := 0
	mark := 0
	t1 := strings.Split(r, " ")
	for k, v := range t1 {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		} else {
			mark++
			if mark == c {
				n = k + 1
				break
			}
		}
	}
	t2 := strings.SplitN(r, " ", n)
	return t2[0], t2[n-1]
}
