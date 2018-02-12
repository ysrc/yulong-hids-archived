// +build windows

package collect

import "github.com/StackExchange/wmi"

type service struct {
	Caption   string // 描述信息
	Name      string // 服务名称
	PathName  string // 服务程序路径
	Started   bool   // 是否已启动
	StartMode string // 启动模式
	StartName string // 启动用户
}

// GetServiceInfo 获取服务列表
func GetServiceInfo() (resultdata []map[string]string) {
	var dst []service
	err := wmi.Query("SELECT * FROM Win32_Service", &dst)
	if err != nil {
		return resultdata
	}
	for _, v := range dst {
		m := make(map[string]string)
		m["name"] = v.Name
		m["pathname"] = v.PathName
		if v.Started {
			m["started"] = "True"
		} else {
			m["started"] = "False"
		}
		m["startmode"] = v.StartMode
		m["startname"] = v.StartName
		m["caption"] = v.Caption
		resultdata = append(resultdata, m)
	}
	return resultdata
}
