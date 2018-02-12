// +build linux

package collect

import (
	"regexp"
	"strings"
	"yulong-hids/agent/common"
)

// GetListening 获取tcp端口监听端口
func GetListening() (resultData []map[string]string) {
	listeningStr := common.Cmdexec("ss -nltp")
	listeningList := strings.Split(listeningStr, "\n")
	if len(listeningList) < 2 {
		return
	}
	for _, info := range listeningList[1 : len(listeningList)-1] {
		if strings.Contains(info, "127.0.0.1") {
			continue
		}
		m := make(map[string]string)
		reg := regexp.MustCompile("\\s+")
		info = reg.ReplaceAllString(strings.TrimSpace(info), " ")
		s := strings.Split(info, " ")
		if len(s) < 6 {
			continue
		}
		m["proto"] = "TCP"
		if strings.Contains(s[3],"::"){
			m["address"] = strings.Replace(s[3], "::", "0.0.0.0", 1)
		}else{
			m["address"] = strings.Replace(s[3], "*", "0.0.0.0", 1)
		}
		b := false
		for _,v:= range resultData{
			if v["address"] == m["address"]{
				b = true
				break
			}
		}
		if b{
			continue
		}
		reg = regexp.MustCompile(`users:\(\("(.*?)",(.*?),.*?\)`)
		r := reg.FindSubmatch([]byte(s[5]))
		if strings.Contains(string(r[2]), "=") {
			m["pid"] = strings.SplitN(string(r[2]), "=", 2)[1]
		} else {
			m["pid"] = string(r[2])
		}
		m["name"] = string(r[1])
		resultData = append(resultData, m)
	}
	return resultData
}
