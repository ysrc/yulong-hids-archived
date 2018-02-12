// +build windows

package collect

import "github.com/StackExchange/wmi"

type userAccount struct {
	Name        string // 用户名
	Description string // 用户描述
	Status      string // 用户状态
}

// GetUser 获取系统用户列表
func GetUser() (resultData []map[string]string) {
	var dst []userAccount
	err := wmi.Query("SELECT * FROM Win32_UserAccount where LocalAccount=TRUE", &dst)
	if err != nil {
		return resultData
	}
	for _, v := range dst {
		m := make(map[string]string)
		m["name"] = v.Name
		m["description"] = v.Description
		m["status"] = v.Status
		resultData = append(resultData, m)
	}
	return resultData
}

/**
func GetUser() []map[string]string {
	var resultData []map[string]string
	k, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SAM\\SAM\\Domains\\Account\\Users\\Names", registry.ALL_ACCESS)
	userList, _ := k.ReadSubKeyNames(0)
	for _, user := range userList {
		m := map[string]string{"name": user}
		//m["name"] = user
		resultData = append(resultData, m)
	}
	return resultData
}
**/
