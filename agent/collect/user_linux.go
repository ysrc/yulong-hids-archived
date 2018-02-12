// +build linux

package collect

import (
	"io/ioutil"
	"strings"
)

// GetUser 获取系统用户列表
func GetUser() (resultData []map[string]string) {
	dat, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return resultData
	}
	userList := strings.Split(string(dat), "\n")
	if len(userList) < 2 {
		return
	}
	for _, info := range userList[0 : len(userList)-2] {
		if strings.Contains(info, "/nologin") {
			continue
		}
		s := strings.SplitN(info, ":", 2)
		m := map[string]string{"name": s[0], "description": s[1]}
		resultData = append(resultData, m)
	}
	return resultData
}
