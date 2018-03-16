package install

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"crypto/md5"
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
	var err error
	url := fmt.Sprintf("%s://%s/json/download?system=%s&platform=%s&type=agent&action=download", common.Proto, ip, runtime.GOOS, arch)

	// Agent 下载检查和重试, 重试三次，功能性考虑
	times := 3
	for {
		err = downFile(url, agentPath)
		// 检查文件hash是否匹配
		if err == nil {
			mstr, _ := FileMD5String(agentPath)
			log.Println("Agent file MD5:", mstr)
			if CheckAgentHash(mstr, ip, arch) {
				log.Println("Agent download finished, hash check passed")
				return nil
			} else {
				log.Println("Agent is broken, retry the downloader again")
			}
		}
		if times--; times == 0 {
			break
		}
	}

	return errors.New("Agent Download Error")
}

// FileMD5String 获取文件MD5
func FileMD5String(filePath string) (MD5String string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5h := md5.New()
	io.Copy(md5h, file)
	return fmt.Sprintf("%x", md5h.Sum([]byte(""))), nil
}


// CheckAgentHash 检查Agent的哈希值是否匹配
func CheckAgentHash(fileHash string, ip string, arch string) (is bool) {
	checkURL := fmt.Sprintf(
		"%s://%s/json/download?hash=%x&system=%s&platform=%s&action=check&type=agent",
		common.Proto, ip, fileHash, runtime.GOOS, arch,
	)
	res, err := common.HTTPClient.Get(checkURL)
	if err != nil {
		return false
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	return "1" == string(result)
}