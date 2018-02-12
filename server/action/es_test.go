package action

import (
	"fmt"
	"testing"
	"time"
	"yulong-hids/server/models"
)

func TestResultSave(t *testing.T) {
	datainfo := models.DataInfo{
		"22.22.22.22",
		"loginlog",
		"windows",
		[]map[string]string{{"user": "test", "time": "2017-08-10 9:59:43"}},
		time.Now(),
	}
	fmt.Println(ResultSave(datainfo))
}

func TestComputerInfoSave(t *testing.T) {
	fmt.Println(time.Unix(1502330143, 0).String())
	loc, _ := time.LoadLocation("Local")
	time, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-08-10 9:55:43", loc)
	//time, _ := time.Parse("2006-01-02 15:04:05", "2017-08-10 9:55:43")
	fmt.Println(time.String(), 1)
}

func TestGetAgentConfig(t *testing.T) {
	fmt.Println(GetAgentConfig("10.100.173.38"))
	t.Error(1)
}
