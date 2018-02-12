// +build linux

package collect

type service struct {
	Caption   string
	Name      string
	PathName  string
	Started   bool
	StartMode string
	StartName string
}

// GetServiceInfo 获取服务列表
func GetServiceInfo() []map[string]string {
	var resultdata []map[string]string
	return resultdata
}
