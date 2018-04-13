package routers

import (
	"wxwork/controllers"
	"github.com/astaxie/beego"
	"wxwork/controllers/api/wxwork"
	"wxwork/controllers/wxwork"
)

func init() {

	beego.Router("/api/wxwork/apps/callback", &controllers_api_wxwork.ApiWxAppsController{},"*:Callback")
	beego.Router("/api/wxwork/apps/:id/callback", &controllers_api_wxwork.ApiWxAppsController{},"*:Callback")
	beego.Router("/api/wxwork/apps/:id/syn_org", &controllers_api_wxwork.ApiWxAppsController{},"get:SynOrg")

	beego.Router("/api/wxwork/orgs/callback", &controllers_api_wxwork.ApiWxOrgController{},"*:Callback")
	beego.Router("/api/wxwork/orgs/:id/callback", &controllers_api_wxwork.ApiWxOrgController{},"*:Callback")


	beego.Router("/wxwork", &controllers_wxwork.DashboardController{},"get:Index")
	beego.Router("/wxwork/home", &controllers_wxwork.DashboardController{},"get:Home")
	beego.Router("/wxwork/:app_id/home", &controllers_wxwork.DashboardController{},"get:Home")

    beego.Router("/", &controllers.MainController{})
}
