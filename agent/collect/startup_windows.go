// +build windows

package collect

import (
	"yulong-hids/agent/common"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows/registry"
)

type startup struct {
	Caption  string // 描述信息
	Command  string // 执行的程序、命令
	Location string // 开机启动来源
	User     string // 启动用户
}

// GetStartup 获取开机启动项
func GetStartup() (resultData []map[string]string) {

	// 通过wmi获取注册表与启动目录方式开机启动列表
	getStartupInWMI(&resultData)

	// 获取ActiveX方式开机启动列表
	getActiveXStatup(&resultData)

	// 获取启动目录方式开机启动列表
	// getStatupPath(&resultData)
	return resultData
}
func getStartupInWMI(resultData *[]map[string]string) {
	var dst []startup
	err := wmi.Query("SELECT * FROM Win32_StartupCommand", &dst)
	if err != nil {
		return
	}
	for _, v := range dst {
		m := make(map[string]string)
		m["name"] = v.Caption
		m["command"] = v.Command
		m["location"] = v.Location
		m["user"] = v.User
		*resultData = append(*resultData, m)
	}
}
func getActiveXStatup(resultData *[]map[string]string) {
	// 获取启动列表
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Active Setup\\Installed Components", registry.ALL_ACCESS|registry.WOW64_64KEY)
	if err != nil {
		return
	}
	uplist, err := k.ReadSubKeyNames(0)
	if err != nil {
		return
	}
	k.Close()
	// 获取已启动列表(已启动的不会再运行)
	k, err = registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Active Setup\\Installed Components", registry.ALL_ACCESS|registry.WOW64_64KEY)
	if err != nil {
		return
	}
	overlist, _ := k.ReadSubKeyNames(0)
	k.Close()
	for _, key := range uplist {
		k2, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Active Setup\\Installed Components\\"+key, registry.ALL_ACCESS|registry.WOW64_64KEY)
		if !common.InArray(overlist, key, false) {
			if command, _, err := k2.GetStringValue("StubPath"); err == nil {
				name, _, _ := k2.GetStringValue("")
				m := map[string]string{"name": name, "command": command, "location": "Active Setup", "user": "Public"}
				*resultData = append(*resultData, m)
			}
		}
		k2.Close()
	}
}

/**
func getRegRun(resultData *[]map[string]string) {
	k, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	uplist, _ := k.ReadValueNames(-1)
	for _, upname := range uplist {
		data, _, _ := k.GetStringValue(upname)
		m := map[string]string{"name": upname, "data": data, "source": "registry"}
		*resultData = append(*resultData, m)
	}
	k, _ = registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	uplist, _ = k.ReadValueNames(-1)
	for _, upname := range uplist {
		data, _, _ := k.GetStringValue(upname)
		m := map[string]string{"name": upname, "data": data, "source": "registry"}
		*resultData = append(*resultData, m)
	}
	if runtime.GOARCH == "amd64" {
		k, _ = registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
		uplist, _ = k.ReadValueNames(-1)
		for _, upname := range uplist {
			data, _, _ := k.GetStringValue(upname)
			//fmt.Println(upname, data)
			m := map[string]string{"name": upname, "data": data, "source": "registry"}
			*resultData = append(*resultData, m)
		}
	}
	k.Close()
}


func getStatupPath(resultData *[]map[string]string) {
	var pathlist []string
	k, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders", registry.ALL_ACCESS)
	data, _, err := k.GetStringValue("Common Startup")
	if err == nil {
		pathlist = append(pathlist, data+`\`)
	}
	k, _ = registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders", registry.ALL_ACCESS)
	data, _, err = k.GetStringValue("Startup")
	if err == nil {
		pathlist = append(pathlist, data+`\`)
	}
	k.Close()

	for _, path := range pathlist {
		dirList, _ := ioutil.ReadDir(path)
		for _, file := range dirList {
			if file.Name() != "desktop.ini" {
				var m map[string]string
				if strings.HasSuffix(file.Name(), ".lnk") {
					dat, _ := ioutil.ReadFile(path + file.Name())
					filepath := strings.Split(strings.Split(string(dat), "System\x00")[1], "\x00")[0]
					m = map[string]string{"name": file.Name(), "data": filepath, "source": "startup"}
				} else {
					m = map[string]string{"name": file.Name(), "data": path + file.Name(), "source": "startup"}
				}
				*resultData = append(*resultData, m)
			}
		}
	}
}
**/
