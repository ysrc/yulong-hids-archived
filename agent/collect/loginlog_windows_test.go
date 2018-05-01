package collect

import (
	"regexp"
	"testing"
)

var re = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

// TestGetLoginLog GetLoginLog的测试模块
/*

command: go test -timeout 30s yulong-hids/agent/collect -run ^TestGetLoginLog$

Fail output in issue: https://github.com/ysrc/yulong-hids/issues/29 :

--- FAIL: TestGetLoginLog (0.06s)
        loginlog_windows_test.go:23: Remote is "65014", want ip address string.
        loginlog_windows_test.go:23: Remote is "59844", want ip address string.
        loginlog_windows_test.go:23: Remote is "59991", want ip address string.
        loginlog_windows_test.go:23: Remote is "60125", want ip address string.
        loginlog_windows_test.go:23: Remote is "63651", want ip address string.
        loginlog_windows_test.go:23: Remote is "54666", want ip address string.
        loginlog_windows_test.go:23: Remote is "64662", want ip address string.
        loginlog_windows_test.go:23: Remote is "64820", want ip address string.
        loginlog_windows_test.go:23: Remote is "55321", want ip address string.
        loginlog_windows_test.go:23: Remote is "51392", want ip address string.
        loginlog_windows_test.go:23: Remote is "53599", want ip address string.
        loginlog_windows_test.go:23: Remote is "50233", want ip address string.
        loginlog_windows_test.go:23: Remote is "52694", want ip address string.
        loginlog_windows_test.go:23: Remote is "63787", want ip address string.
FAIL
FAIL    yulong-hids/agent/collect       0.094s

*/
func TestGetLoginLog(t *testing.T) {
	loginInfoLst := GetLoginLog()
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
}
