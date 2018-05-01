package collect

import (
	"fmt"
	"regexp"
	"testing"
)

var re = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

// TestGetLoginLog GetLoginLog的测试模块
func TestGetLoginLog(t *testing.T) {
	loginInfoLst := GetLoginLog()
	fmt.Println(loginInfoLst)
	for _, loginInfo := range loginInfoLst {
		// status 只能是 true 或 false
		status := loginInfo["status"]
		if !(status == "true" || status == "false") {
			t.Errorf("Status is %q, want bool string.", status)
		}
		// remote 应该是 IP
		remote := loginInfo["remote"]
		if !(re.MatchString(remote)) {
			t.Errorf("Remote is %q, want ip address string.", remote)
		}
	}
	fmt.Println("TestGetLoginLog all test passed!!!")
}
