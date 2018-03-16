// +build windows

package install

import (
	"log"
	"strings"
	"yulong-hids/daemon/common"
)

// Agent 下载安装agent
func Agent(ip string, installPath string, arch string) error {
	// 下载agent.exe
	log.Println("Download Agent")
	err := DownAgent(ip, installPath+"agent.exe", arch)
	if err != nil {
		return err
	}
	// 拷贝自身到安装目录
	log.Println("Copy the daemon to the installation directory")
	err = copyMe(installPath)
	if err != nil {
		return err
	}
	// 安装daemon为服务
	cmd := installPath + "daemon.exe -register -netloc " + ip
	out, err := common.CmdExec(cmd)
	if err != nil {
		return err
	}
	// 启动服务
	log.Println("Start the service")
	cmd = "net start yulong-hids"
	out, err = common.CmdExec(cmd)
	if err == nil && strings.Contains(out, "yulong-hids") {
		log.Println("Start service successfully")
		return nil
	}
	return err
}
