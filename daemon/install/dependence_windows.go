// +build windows

package install

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"yulong-hids/daemon/common"

	"github.com/StackExchange/wmi"
)

type service struct {
	Name string
}

func hasService(name string) bool {
	var dst []service
	err := wmi.Query(fmt.Sprintf(`SELECT * FROM Win32_Service where Name = "%s"`, name), &dst)
	if err == nil && len(dst) == 1 {
		return true
	}
	return false
}

// Dependency 下载->释放->安装依赖文件和服务
func Dependency(ip string, installPath string, arch string) error {
	url := fmt.Sprintf("%s://%s/json/download?system=windows&platform=%s&type=data&action=download", common.Proto, ip, arch)
	pcappath := installPath + "data.zip"
	log.Println("Download dependent environment package")
	err := downFile(url, pcappath)
	if err != nil {
		return err
	}
	rc, err := zip.OpenReader(pcappath)
	if err != nil {
		return err
	}
	defer rc.Close()
	for _, _file := range rc.File {
		f, err := _file.Open()
		if err != nil {
			return err
		}
		desfile, err := os.OpenFile(installPath+_file.Name, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err == nil {
			io.CopyN(desfile, f, int64(_file.UncompressedSize64))
			desfile.Close()
		} else {
			return err
		}
	}
	if !hasService("npf") {
		log.Println("Install npf service")
		common.CmdExec(fmt.Sprintf("sc create npf binPath= %s/npf.sys type= kernel start= auto error= normal", installPath))
		common.CmdExec("net start npf")
	}
	if !hasService("pro") {
		log.Println("Install pro service")
		common.CmdExec(fmt.Sprintf("sc create pro binPath= %s/pro.sys type= kernel start= auto error= normal", installPath))
		common.CmdExec("net start pro")
	}
	return nil
}
