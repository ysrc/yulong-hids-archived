package task

import (
	"log"
	"os"
	"runtime"
	"yulong-hids/daemon/common"
)

// UnInstallALL 卸载
func UnInstallALL() {
	if runtime.GOOS == "windows" {
		common.CmdExec("net stop pro")
		// common.CmdExec("net stop npf")
		common.CmdExec("sc delete pro")
		// common.CmdExec("sc delete npf")
	} else {
		common.CmdExec("rmmod syshook_execve")
	}
	common.KillAgent()
	if err := common.Service.Uninstall(); err != nil {
		log.Println("Uninstall yulong-hids error:", err.Error())
	}
	log.Println("Uninstall completed")
	os.Exit(1)
}
