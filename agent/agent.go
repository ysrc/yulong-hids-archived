package main

import (
	"fmt"
	"log"
	"os"
	"yulong-hids/agent/client"
	"runtime"
	"yulong-hids/daemon/common"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: agent[.exe] ServerIP [debug]")
		fmt.Println("Example: agent 8.8.8.8 debug")
		return
	}
	if runtime.GOOS == "linux" {
		out, _ := common.CmdExec(fmt.Sprintf("lsmod|grep syshook_execve"))
		if out == "" {
			common.CmdExec(fmt.Sprintf("insmod %s/syshook_execve.ko", common.InstallPath))
		}
	}
	var agent client.Agent
	agent.ServerNetLoc = os.Args[1]
	if len(os.Args) == 3 && os.Args[2] == "debug" {
		log.Println("DEBUG MODE")
		agent.IsDebug = true
	}
	agent.Run()
}
