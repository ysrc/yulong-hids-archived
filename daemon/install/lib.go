package install

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"yulong-hids/daemon/common"
)

//下载文件
func downFile(url string, svaepath string) error {
	request, _ := http.NewRequest("GET", url, nil)
	request.Close = true
	if res, err := common.HTTPClient.Do(request); err == nil {
		defer res.Body.Close()
		file, err := os.Create(svaepath)
		if err != nil {
			return err
		}
		io.Copy(file, res.Body)
		file.Close()
		if runtime.GOOS == "linux" {
			os.Chmod(svaepath, 0750)
		}
		fileInfo, err := os.Stat(svaepath)
		// log.Println(res.ContentLength, fileInfo.Size())
		if err != nil || fileInfo.Size() != res.ContentLength {
			log.Println("File download error:", err.Error())
			return errors.New("downfile error")
		}
	} else {
		return err
	}
	return nil
}

//复制自身
func copyMe(installPath string) (err error) {
	var dstName string
	if runtime.GOOS == "windows" {
		dstName = installPath + "daemon.exe"
	} else {
		dstName = installPath + "daemon"
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	mepath, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	if mepath == dstName {
		return nil
	}
	src, err := os.Open(mepath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

// DownAgent 下载agent到指定安装目录
func DownAgent(ip string, agentPath string, arch string) error {
	url := fmt.Sprintf("%s://%s/json/download?system=%s&platform=%s&type=agent&action=download", common.Proto, ip, runtime.GOOS, arch)
	err := downFile(url, agentPath)
	return err
}
