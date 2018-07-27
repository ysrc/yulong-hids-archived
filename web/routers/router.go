package routers

import (
	"yulong-hids/web/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/"+beego.AppConfig.String("ApiVer"),
		beego.NSRouter("/client", &controllers.ClientController{}),
		beego.NSRouter("/download", &controllers.DloadController{}),
		beego.NSRouter("/serverlist", &controllers.AgentApiController{}),
		beego.NSRouter("/myip", &controllers.AgentApiController{}),
		beego.NSRouter("/publickey", &controllers.AgentApiController{}),
		beego.NSRouter("/dbinfo", &controllers.AgentApiController{}),
		beego.NSRouter("/statistics", &controllers.StatisticsController{}),
		beego.NSRouter("/file", &controllers.FileController{}, "post:Upload"),
		beego.NSRouter("/analyze", &controllers.AnalyzeController{}, "post:Post;get:Get"),
		beego.NSRouter("/config", &controllers.ConfigController{}, "get:Get;post:Edit;delete:Del;put:Add"),
		beego.NSRouter("/info/:ip", &controllers.InfoController{}, "get:GetInfoByIp"),
		beego.NSRouter("/monitor/:ip/:type/:start", &controllers.MonitorController{}, "get:GetTwenty"),
		beego.NSRouter("/monitor/:ip", &controllers.MonitorController{}, "get:GetAllType"),
		beego.NSRouter("/notice", &controllers.NoticeController{}, "get:Get;post:ChangeStatus;delete:Delete"),
		beego.NSRouter("/tasks", &controllers.TaskController{}, "get:Get;post:Post"),
		beego.NSRouter("/rules", &controllers.RuleController{}, "get:Get;post:Post"),
		beego.NSRouter("/logout", &controllers.LogoutController{}, "post:Post"),
	)
	beego.AddNamespace(ns)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login/", &controllers.LoginController{})
	beego.Router("/install/", &controllers.InstallController{})

}
