// +build windows

package collect

import (
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"yulong-hids/agent/common"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
func getWebPath(webCommand string) ([]string, error) {
	var pathList []string
	var iis6ConfigPath string
	var iis7ConfigPath string
	SystemDrive := os.Getenv("SystemDrive")
	if runtime.GOARCH == "386" {
		iis6ConfigPath = SystemDrive + `\WINDOWS\SysNative\inetsrv\MetaBase.xml`
		iis7ConfigPath = SystemDrive + `\Windows\SysNative\inetsrv\config\applicationHost.config`
	} else {
		iis6ConfigPath = SystemDrive + `\WINDOWS\System32\inetsrv\MetaBase.xml`
		iis7ConfigPath = SystemDrive + `\Windows\System32\inetsrv\config\applicationHost.config`
	}
	//IIS 6
	if pathExists(iis6ConfigPath) {
		dat, err := ioutil.ReadFile(iis6ConfigPath)
		if err != nil {
			return pathList, err
		}
		reg := regexp.MustCompile(`Path="(.*?)"\s`)
		pathM := reg.FindAllSubmatch([]byte(dat), -1)
		for _, info := range pathM {
			if common.InArray(pathList, string(info[1]), false) {
				pathList = append(pathList, string(info[1]))
			}
		}
	}
	//IIS 7
	if pathExists(iis7ConfigPath) {
		dat, err := ioutil.ReadFile(iis7ConfigPath)
		if err != nil {
			return pathList, err
		}
		reg := regexp.MustCompile(`physicalPath="(.*?)"`)
		pathM := reg.FindAllSubmatch([]byte(dat), -1)
		for _, info := range pathM {
			if !common.InArray(pathList, string(info[1]), false) {
				pathList = append(pathList, strings.Replace(string(info[1]), "%SystemDrive%", SystemDrive, -1))
			}
		}
	}
	return pathList, nil
}
