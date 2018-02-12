// +build linux
package monitor

/*
#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>

static  struct   stat buf;

int  get_id(char * filename)
{
    stat(filename, &buf);
    int uid=buf.st_uid;
    return uid;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"yulong-hids/agent/common"

	"github.com/go-fsnotify/fsnotify"
)

func getFileUser(path string) (string, error) {
	uidStr := fmt.Sprintf("%d", C.get_id(C.CString(path)))
	dat, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return "", err
	}
	userList := strings.Split(string(dat), "\n")
	for _, info := range userList[0 : len(userList)-2] {
		// fmt.Println(info)
		s := strings.SplitN(info, ":", -1)
		if len(s) >= 3 && s[2] == uidStr {
			// fmt.Println(s[0])
			return s[0], nil
		}
	}
	return "", errors.New("error")
}

// StartFileMonitor 开始文件行为监控
func StartFileMonitor(resultChan chan map[string]string) {
	log.Println("StartFileMonitor")
	var pathList []string
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	for _, path := range common.Config.MonitorPath {
		if path == "%web%" {
			iterationWatcher(common.ServerInfo.Path, watcher, pathList)
			continue
		}
		if strings.HasPrefix(path, "/") {
			pathList = append(pathList, path)
			if strings.HasSuffix(path, "*") {
				iterationWatcher([]string{strings.Replace(path, "*", "", 1)}, watcher, pathList)
			} else {
				watcher.Add(path)
			}
		}
	}
	var resultdata map[string]string
	for {
		select {
		case event := <-watcher.Events:
			resultdata = make(map[string]string)
			if common.InArray(filter.File, strings.ToLower(event.Name), false) ||
				common.InArray(pathList, strings.ToLower(event.Name), false) ||
				common.InArray(common.Config.Filter.File, strings.ToLower(event.Name), true) {
				continue
			}
			if len(event.Name) == 0 {
				continue
			}
			resultdata["source"] = "file"
			resultdata["action"] = event.Op.String()
			resultdata["path"] = event.Name
			resultdata["hash"] = ""
			resultdata["user"] = ""
			f, err := os.Stat(event.Name)
			if err == nil && !f.IsDir() {
				if f.Size() <= fileSize {
					if hash, err := getFileMD5(event.Name); err == nil {
						resultdata["hash"] = hash
						if common.InArray(common.Config.Filter.File, strings.ToLower(hash), false) {
							continue
						}
					}
				}
				if user, err := getFileUser(event.Name); err == nil {
					resultdata["user"] = user
				}
			}
			if isFileWhite(resultdata) {
				continue
			}
			resultChan <- resultdata
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
