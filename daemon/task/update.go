package task

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
	"yulong-hids/daemon/common"
	"yulong-hids/daemon/install"
)

func agentUpdate(ip string, installPath string, arch string) (bool, error) {
	var err error
	var agentFilePath string
	if runtime.GOOS == "windows" {
		agentFilePath = installPath + "agent.exe"
	} else {
		agentFilePath = installPath + "agent"
	}
	file, err := os.Open(agentFilePath)
	if err == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		file.Close()
		agentMd5 := md5h.Sum([]byte(""))
		checkURL := fmt.Sprintf("%s://%s/json/download?hash=%x&system=%s&platform=%s&action=check&type=agent", common.Proto, ip, agentMd5, runtime.GOOS, arch)
		log.Println(checkURL)
		res, err := common.HTTPClient.Get(checkURL)
		if err != nil {
			return false, err
		}
		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		if string(result) == "1" {
			common.M.Lock()
			defer common.M.Unlock()
			log.Println("Updated Agent")
			common.KillAgent()
			time.Sleep(time.Second)
			if err = os.Remove(agentFilePath); err == nil {
				if err = install.DownAgent(ip, agentFilePath, arch); err == nil {
					log.Println("Download replacement success")
					return true, nil
				}
			} else {
				return false, err
			}
		}
	}
	return false, err
}
