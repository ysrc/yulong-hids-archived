// +build linux

package install

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"yulong-hids/daemon/common"
)

func Dependency(ip string, installPath string, arch string) error {
	// _, err := os.Stat("/usr/lib64/libpcap.so.1")
	// if err != nil {
	// 	return err
	// }
	url := fmt.Sprintf("%s://%s/json/download?system=linux&platform=%s&type=data&action=download", common.Proto, ip, arch)
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
	out, err := common.CmdExec("uname -r")
	if err != nil || strings.Count(out, ".") < 2 {
		return errors.New("Get kernel version identification")
	}
	ver := strings.Join(strings.Split(strings.Trim(out, "\n"), ".")[0:3], ".")
	for _, _file := range rc.File {
		if _file.Name == "syshook_"+ver+".ko" {
			f, _ := _file.Open()
			desfile, err := os.OpenFile(installPath+"syshook_execve.ko", os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			io.CopyN(desfile, f, int64(_file.UncompressedSize64))
			desfile.Close()
			log.Println("Use syshook_" + ver)
			out, err = common.CmdExec(fmt.Sprintf("insmod %s/syshook_execve.ko", installPath))
			if err != nil {
				return err
			}
			if !strings.Contains(out, "ERROR") {
				log.Println("Insmod syshook_execve succeeded")
			} else {
				log.Println("Insmod syshook_execve error, command output:", out)
			}
		}
	}
	return nil
}
