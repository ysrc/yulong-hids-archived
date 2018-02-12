package safecheck

import (
	"fmt"
	"testing"
	"time"
	"yulong-hids/server/models"
)

// func init() {
// 	data
// }

func BenchmarkCheck(b *testing.B) {
	// testdata := models.DataInfo{"10.101.20.73", "connection", "windows",
	// 	[]map[string]string{{"ip": "123.125.65.153:80", "local": "10.101.20.73:50596"}}, time.Now()}
	testdata := models.DataInfo{"10.101.20.73", "process", "windows",
		[]map[string]string{{"name": "svchost.exe", "command": "C:\\Windows\\System32\\svchost.exe -k WerSvcGroup|svchost.exe", "parentname": "cmd.exe"}}, time.Now()}
	var c Check
	c.CStatistics = models.DB.C("statistics")
	c.CNoice = models.DB.C("notice")
	for i := 0; i < b.N; i++ {
		c.Info = testdata
		c.Run()
	}
}
func Test_check(t *testing.T) {
	// testdata := models.DataInfo{"10.101.20.73", "connection", "windows",
	// 	[]map[string]string{{"ip": "123.125.65.153:80", "local": "10.101.20.73:50596"}}, time.Now()}
	testdata := models.DataInfo{"10.1.1.1", "process", "windows",
		[]map[string]string{{"name": "mimikatz.exe", "command": `mimikatz.exe "privilege::debug" "sekurlsa::logonPasswords full" exit`, "parentname": "cmd.exe"}}, time.Now()}
	// testdata := models.DataInfo{"10.101.20.73", "file", "windows",
	// 	[]map[string]string{{"action": "WRITE",
	// 		"hash": "74e9048db945f3d075fccbb3aace735d",
	// 		"path": "c:\\sitefile\\portal20170804\\content\\data\\test.aspx",
	// 		"user": "protal"}}, time.Now()}
	//ScanChan <- testdata
	fmt.Println(models.Config)
	c := new(Check)
	c.CStatistics = models.DB.C("statistics")
	c.CNoice = models.DB.C("notice")
	c.Info = testdata
	c.Run()
	t.Error(1)
}

func Test_DownCheckThread(t *testing.T) {
	DownCheckThread()
	t.Error(1)
}
