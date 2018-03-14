package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"yulong-hids/daemon/common"
	"yulong-hids/daemon/install"
	"yulong-hids/daemon/task"

	"github.com/kardianos/service"
)

var (
	ip             *string
	installBool    *bool
	uninstallBool  *bool
	registeredBool *bool
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	go task.WaitThread()
	var agentFilePath string
	if runtime.GOOS == "windows" {
		agentFilePath = common.InstallPath + "agent.exe"
	} else {
		agentFilePath = common.InstallPath + "agent"
	}
	for {
		common.M.Lock()
		log.Println("Start Agent")
		common.Cmd = exec.Command(agentFilePath, common.ServerIP)
		err := common.Cmd.Start()
		common.M.Unlock()
		if err == nil {
			common.AgentStatus = true
			log.Println("Start Agent successful")
			err = common.Cmd.Wait()
			if err != nil {
				common.AgentStatus = false
				log.Println("Agent to exit：", err.Error())
			}
		} else {
			log.Println("Startup Agent failed", err.Error())
		}
		time.Sleep(time.Second * 10)
	}
}

func (p *program) Stop(s service.Service) error {
	common.KillAgent()
	return nil
}

func main() {
	flag.StringVar(&common.ServerIP, "netloc", "", "* WebServer 192.168.1.100:443")
	installBool = flag.Bool("install", false, "Install yulong-hids service")
	uninstallBool = flag.Bool("uninstall", false, "Remove yulong-hids service")
	registeredBool = flag.Bool("register", false, "Registration yulong-hids service")
	flag.Parse()
	svcConfig := &service.Config{
		Name:        "yulong-hids",
		DisplayName: "yulong-hids",
		Description: "集实时监控、异常检测、集中管理为一体的主机安全监测系统",
		Arguments:   []string{"-netloc", common.ServerIP},
	}
	prg := &program{}
	var err error
	common.Service, err = service.New(prg, svcConfig)
	if err != nil {
		log.Println("New a service error:", err.Error())
		return
	}
	if *uninstallBool {
		task.UnInstallALL()
		return
	}
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		return
	}

	// 释放agent
	if *installBool {
		// 依赖环境安装
		if _, err = os.Stat(common.InstallPath); err != nil {
			os.Mkdir(common.InstallPath, 0)
			err = install.Dependency(common.ServerIP, common.InstallPath, common.Arch)
			if err != nil {
				log.Println(err.Error())
			}
		}
		if common.ServerIP == "" {
			flag.PrintDefaults()
			return
		}
		err := install.Agent(common.ServerIP, common.InstallPath, common.Arch)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("Installed")
		return
	}
	// 安装daemon为服务
	if *registeredBool {
		err = common.Service.Install()
		if err != nil {
			log.Println(err.Error())
		} else {
			if err = common.Service.Start(); err != nil {
				log.Println(err.Error())
			} else {
				log.Println("Install as a service", "ok")
			}
		}
		return
	}
	err = common.Service.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
