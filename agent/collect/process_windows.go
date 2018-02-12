// +build windows

package collect

import (
	"fmt"
	"os"
	"time"

	"github.com/StackExchange/wmi"
)

type process struct {
	Name            string
	CommandLine     *string
	ProcessId       uint32
	ParentProcessId uint32
	CreationDate    time.Time
}
type process2003 struct {
	Name            string
	CommandLine     *string
	ProcessId       uint32
	ParentProcessId uint32
}

// GetProcessList 获取当前进程列表，与monitor process同类，故只用于保存显示
func GetProcessList() (resultData []map[string]string) {
	var m map[string]string
	if _, err := os.Stat(os.Getenv("SystemDrive") + `/Users`); err != nil {
		var dst []process2003
		err := wmi.Query("SELECT * FROM Win32_Process", &dst)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(dst) != 0 {
			for _, v := range dst {
				m = make(map[string]string)
				m["name"] = v.Name
				m["pid"] = fmt.Sprintf("%d", v.ProcessId)
				m["ppid"] = fmt.Sprintf("%d", v.ParentProcessId)
				m["command"] = *v.CommandLine
				m["starttime"] = ""
				resultData = append(resultData, m)
			}
		}
	} else {
		var dst []process
		err := wmi.Query("SELECT * FROM Win32_Process", &dst)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(dst) != 0 {
			for _, v := range dst {
				m = make(map[string]string)
				m["name"] = v.Name
				m["pid"] = fmt.Sprintf("%d", v.ProcessId)
				m["ppid"] = fmt.Sprintf("%d", v.ParentProcessId)
				m["command"] = *v.CommandLine
				m["starttime"] = v.CreationDate.String()
				resultData = append(resultData, m)
			}
		}
	}
	return resultData
}
