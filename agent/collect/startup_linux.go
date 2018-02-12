// +build linux

package collect

type Startup struct {
	Caption  string
	Command  string
	Location string
	User     string
}

func GetStartup() []map[string]string {
	var resultData []map[string]string
	return resultData
}
