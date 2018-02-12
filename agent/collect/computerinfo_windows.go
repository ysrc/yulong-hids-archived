// +build windows

package collect

import (
	"fmt"
	"os"
	"yulong-hids/agent/common"

	"golang.org/x/sys/windows/registry"
)

// GetComInfo 获取计算机信息
func GetComInfo() (info common.ComputerInfo) {
	var arch string
	var productName string
	var csdVersion string
	info.IP = common.LocalIP
	info.Hostname, _ = os.Hostname()
	if _, err := os.Stat(os.Getenv("SystemDrive") + `/Windows/SysWOW64/`); err != nil {
		arch = "32"
	} else {
		arch = "64"
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.ALL_ACCESS|registry.WOW64_64KEY)
	if err == nil {
		productName, _, _ = k.GetStringValue("ProductName")
		csdVersion, _, _ = k.GetStringValue("CSDVersion")
		k.Close()
	}
	info.System = fmt.Sprintf("%s %s %s", productName, csdVersion, arch)
	discern(&info)
	return info
}
