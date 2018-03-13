// +build windows

package collect

import (
	"fmt"
	"strconv"
	"time"
	"log"
	"github.com/StackExchange/wmi"
)

// 在 Windows 2003 里有可能无法获得创建时间：CreationDate
// 实际测试过程中, Windows 2008 R2 也有可能无法获取到时间导致错误
type process struct {
	Name            string
	CommandLine     *string
	ProcessId       uint32
	ParentProcessId uint32
	CreationDate    *string
}


// GetProcessList 获取当前进程列表，与monitor process同类，故只用于保存显示
func GetProcessList() (resultData []map[string]string) {
	var dst []process
	err := wmi.Query("SELECT * FROM Win32_Process", &dst)
	if err != nil {
		log.Println("Windows get process info wmi Query error:",err)
		return
	}
	if len(dst) != 0 {
		for _, v := range dst {
			m := make(map[string]string)
			m["name"] = v.Name
			m["pid"] = fmt.Sprintf("%d", v.ProcessId)
			m["ppid"] = fmt.Sprintf("%d", v.ParentProcessId)
			m["command"] = *v.CommandLine
			m["starttime"] = parseTime(*v.CreationDate)
			log.Println("Process data:", m)
			resultData = append(resultData, m)
		}
	}
	return resultData
}

// parseTime 格式化时间字符串
func parseTime(val string) string {
	if val == "" {
		return ""
	}
	if len(val) == 25 {
		// 数据格式例如 20180312092136.238259+480
		mins, err := strconv.Atoi(val[22:])
		if err != nil {
			return ""
		}
		val = val[:22] + fmt.Sprintf("%02d%02d", mins/60, mins%60)
	}
	t, err := time.Parse("20060102150405.000000-0700", val)
	if err != nil {
		return ""
	}
	return t.String()
}
