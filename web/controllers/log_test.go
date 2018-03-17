package controllers

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
)

func TestLogLevel(t *testing.T) {
	fmt.Println("beego.LevelAlert : ", beego.LevelAlert)
	fmt.Println("beego.LevelCritical : ", beego.LevelCritical)
	fmt.Println("beego.LevelError : ", beego.LevelError)
	fmt.Println("beego.LevelWarning : ", beego.LevelWarning)
	fmt.Println("beego.LevelNotice : ", beego.LevelNotice)
	fmt.Println("beego.LevelInformational : ", beego.LevelInformational)
	fmt.Println("beego.LevelDebug : ", beego.LevelDebug)
}
