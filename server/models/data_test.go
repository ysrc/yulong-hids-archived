package models

import "fmt"
import "testing"

func TestQueryLogLastTime(t *testing.T) {
	fmt.Println(QueryLogLastTime("10.101.20.70"))
}

func Test_esCheck(t *testing.T) {
	esCheck()
	t.Error(1)
}
