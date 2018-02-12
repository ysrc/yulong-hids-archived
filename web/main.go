package main

import (
	_ "yulong-hids/web/routers"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"
)

var (
	logFile         string
	logConfJSON     string
	defaultLogLevel int
)

func init() {
	defaultLogLevel = beego.LevelInformational
	logFile = beego.AppConfig.String("logfile")
	logConfJSON = `{"filename":"` + logFile + `"}`
}

func main() {
	beego.SetLogger("file", logConfJSON)

	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = settings.SessionGCMaxLifetime
	settings.ProjectPath = utils.GetCwd()
	settings.FilePath = utils.DloadFilePath(settings.ProjectPath)

	// set loglevel
	if beego.AppConfig.String("runmode") == "dev" {
		beego.SetLevel(beego.LevelDebug)
	} else if level, err := beego.AppConfig.Int("loglevel"); err == nil {
		beego.SetLevel(level)
	} else {
		beego.SetLevel(defaultLogLevel)
	}

	// add /tests to https://domain/tests as static path in develop mode
	if utils.IsDevMode() {
		beego.SetStaticPath("/tests", "tests")
	}

	beego.Run()
}
