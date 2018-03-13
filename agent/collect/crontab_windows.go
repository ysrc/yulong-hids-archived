// +build windows

package collect

import (
	"encoding/xml"
	"log"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/axgle/mahonia"
)

type task struct {
	RegistrationInfo struct {
		Description string
	}
	Actions struct {
		Exec struct {
			Command   string
			Arguments string
		}
	}
	Triggers struct {
		CalendarTrigger struct {
			StartBoundary string
		}
	}
	Principals struct {
		Principal struct {
			UserId string
		}
	}
}

// GetCrontab 获取计划任务
func GetCrontab() (resultData []map[string]string) {
	//系统计划任务
	var taskPath string
	if runtime.GOARCH == "386" {
		taskPath = `C:\Windows\SysNative\Tasks\`
	} else {
		taskPath = `C:\Windows\System32\Tasks\`
	}
	dir, err := ioutil.ReadDir(taskPath)
	if err != nil {
		return resultData
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		dat, err := ioutil.ReadFile(taskPath + f.Name())
		if err != nil {
			continue
		}
		v := task{}
		dec := mahonia.NewDecoder("utf-16")
		data := dec.ConvertString(string(dat))
		err = xml.Unmarshal([]byte(strings.Replace(data, "UTF-16", "UTF-8", 1)), &v)
		if err != nil {
			log.Println("Windows crontab info xml Unmarshal error: ", err.Error())
			continue
		}
		m := make(map[string]string)
		m["name"] = f.Name()
		m["command"] = v.Actions.Exec.Command
		m["arg"] = v.Actions.Exec.Arguments
		m["user"] = v.Principals.Principal.UserId
		m["rule"] = v.Triggers.CalendarTrigger.StartBoundary
		m["description"] = v.RegistrationInfo.Description
		resultData = append(resultData, m)
	}
	return resultData
}
