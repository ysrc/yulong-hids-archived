package utils

import (
	"github.com/astaxie/beego/logs"
)

// Logger log
var Logger = logs.NewLogger()

// Loginit as name
func Loginit(logfile string) {
	Logger.SetLogger("console")
	Logger.SetLogger(logs.AdapterFile, `{"filename":"info.log","daily":false,"maxdays":365,"level":3}`)
	Logger.EnableFuncCallDepth(true)
}
