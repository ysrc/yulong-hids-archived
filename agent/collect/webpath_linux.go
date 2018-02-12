// +build linux

package collect

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"yulong-hids/agent/common"
)

// 待支持其他webserver
func getWebPath(webCommand string) ([]string, error) {
	var pathList []string
	if ok, _ := regexp.MatchString(`httpd|apache`, webCommand); ok {
		out := common.Cmdexec("apachectl -V")
		if !strings.Contains(string(out), "SERVER_CONFIG_FILE") {
			return pathList, errors.New("Get ConfigFilePath Error!")
		}
		reg2 := regexp.MustCompile(`HTTPD_ROOT="(.*?)"`)
		reg := regexp.MustCompile(`SERVER_CONFIG_FILE="(.*?)"`)
		configFilePath := reg2.FindStringSubmatch(string(out))[1] + "/" + reg.FindStringSubmatch(string(out))[1]
		if configFilePath != "/" {
			dat, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				return pathList, err
			}
			reg = regexp.MustCompile(`<Directory "(.*?)">`)
			pathM := reg.FindAllSubmatch([]byte(dat), -1)
			for _, info := range pathM {
				pathList = append(pathList, string(info[1]))
			}
		}
	} else if strings.Contains(webCommand, "nginx") {
		out := common.Cmdexec("nginx -V")
		regex, _ := regexp.Compile(`\-\-conf\-path\=(.*?)[ |$]`)
		result := regex.FindStringSubmatch(out)
		if len(result) >= 2 {
			configFilePath := result[1]
			dat, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				return pathList, err
			}
			pathRegex, _ := regexp.Compile(`root (.*?)\;`)
			pathResult := pathRegex.FindStringSubmatch(string(dat))
			if len(pathResult) >= 2 {
				pathList = append(pathList, pathResult[1])
			}
			sitePath := filepath.Dir(configFilePath) + "/sites-available/"
			dirList, _ := ioutil.ReadDir(sitePath)
			for _, v := range dirList {
				if v.IsDir() {
					continue
				}
				dat, err := ioutil.ReadFile(sitePath + v.Name())
				if err != nil {
					continue
				}
				pathResult = pathRegex.FindStringSubmatch(string(dat))
				if len(pathResult) >= 2 {
					pathList = append(pathList, pathResult[1])
				}
			}
		}
	}
	return pathList, nil
}
