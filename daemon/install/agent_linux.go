// +build linux
package install

import (
	"log"
	"os"
	"yulong-hids/daemon/common"
)

func Agent(ip string, installPath string, arch string) error {
	// 下载agent
	log.Println("Download Agent")
	err := DownAgent(ip, installPath+"agent", arch)
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
	os.Chmod(installPath+"daemon", 0750)
	cmd := installPath + "daemon -register -netloc " + ip
	out, err := common.CmdExec(cmd)
	if err != nil {
		return err
	}
	//启动服务
	log.Println("Start the service")
	cmd = "systemctl start yulong-hids"
	out, err = common.CmdExec(cmd)
	if err == nil && len(out) == 0 {
		log.Println("Start service successfully")
	}
	return nil
}
