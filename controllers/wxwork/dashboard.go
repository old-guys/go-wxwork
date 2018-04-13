package controllers_wxwork

import (
	"wxwork/services/wxwork"
	"wxwork/models"
	"fmt"
	"strings"
	"net/url"
	"strconv"
	"github.com/astaxie/beego"
)

type DashboardController struct {
	BaseController
	app models.WxworkApp
	wxOrg models.WxworkOrg
	jsapiService services_wxwork.JsapiService
	user models.User
}

func (c *DashboardController) NestPrepare() {
	fmt.Println("Request URL:", c.Ctx.Request.URL)

	actionNames := strings.Join([]string{"Home", "Mobile"}, "#")
	if strings.Contains(actionNames, c.ActionName) { c.findApp() }

	actionNames = strings.Join([]string{"Home", "Mobile"}, "#")
	if strings.Contains(actionNames, c.ActionName) { c.findWxOrg() }

	actionNames = strings.Join([]string{"Home", "Mobile"}, "#")
	if strings.Contains(actionNames, c.ActionName) { c.findJsapiService() }
}

func (c *DashboardController) Index() {
	fmt.Println(c.GetSession("xxxx"))
	c.SetSession("xxxx", "1")

	redirectUri := "http%3A%2F%2Ftest.work.99zmall.com%2Fwxwork%3Fagent_id%3D1000005%26app_id%3D1%26corp_id%3Dwwea672916e0a3a7c4%26mobile_redirect_uri%3D%252Fwxwork%26pc_redirect_uri%3D%252Fwxwork"
	url1 := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=wwea672916e0a3a7c4&agentid=1000005&redirect_uri=" + redirectUri

	c.Data["redirectUri"] = redirectUri
	c.Data["url1"] = url1
	c.TplName = "wxwork/dashboard/index.html"
}

func (c *DashboardController) Home() {
	code := c.Ctx.Input.Query("code")
	userid := c.Ctx.Input.Query("userid")

	if code != "" && !c.MicroMessengerBrowser() {
		c.Redirect("/", 301)
	}

	errCode := ""
	if code == "" && userid == "" {
		c.redirectToAuthorize()
	} else {
		if c.wxOrg.Id == 0 {
			errCode = "org_inactive"
			goto LOOP
		}

		wxUser := models.WxworkUser{}
		if code != "" {
			wxUser = c.jsapiService.UserInfoByCode(code)
		}

		if wxUser.Id == 0 && userid != "" {
			wxUser = c.jsapiService.UserInfoByUserid(userid, map[string]interface{}{})
		}

		if wxUser.Id == 0 { errCode = "user_blank" }
		c.user = wxUser.User()

		if errCode == "" && (c.user.Id == 0 || c.user.Status == 0) {
			errCode = "user_hide"
		}

		if errCode != "" { goto LOOP }

		if c.isMobile() { c.redirectToMobile() }
		if c.isPc() { c.redirectToPc() }
	}

	LOOP:

	c.Data["wxJsapiConfig"] = c.jsapiService.GenerateConfig(c.Ctx.Request.URL.String())
	c.TplName = "wxwork/dashboard/home.html"
}

func (c *DashboardController) redirectToAuthorize() {

	u, _ := url.Parse(beego.AppConfig.String("wxwork_open_host"))
	params := url.Values{}

	params.Add("appid", c.wxOrg.CorpId)
	params.Add("redirect_uri", c.homeRedirectUri())
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_userinfo")
	params.Add("agentid", strconv.Itoa(c.app.AgentID))
	params.Add("state", "STATE")

	u.Path = "/connect/oauth2/authorize"
	u.Fragment = "wechat_redirect"
	u.RawQuery = params.Encode()

	fmt.Println("authorize url:", u.String())
	c.Redirect(u.String(), 301)
}

func (c *DashboardController) redirectToMobile() {
	c.Ctx.Output.Body([]byte("wxwork/mobile"))
}

func (c *DashboardController) redirectToPc() {
	if c.user.Id == 0 {
		c.Ctx.Output.Body([]byte("获取当前登录用户信息失败"))
	}

	// 用户登录sign_out(c.user), sign_in(c.user)

	location := "/"
	redirectUrl := c.Ctx.Input.Query("pc_redirect_uri")
	uri, _ := url.Parse(redirectUrl)

	if redirectUrl != "" && (uri.Host == "" || uri.Host == c.Ctx.Request.Host) {
		location = uri.Path
	}

	fmt.Println(location)
	//c.Redirect(location, 301)
}

func (c *DashboardController) findApp() {
	appId := c.Ctx.Input.Query("app_id")
	if appId == "" { appId = c.Ctx.Input.Param(":app_id")}
	if appId == "" { return }

	models.WxworkAppAr.Where("id = ?", appId).First(&c.app)
}

func (c *DashboardController) findWxOrg() {
	models.WxworkAppAr.Model(&c.app).Related(&c.wxOrg)
}

func (c *DashboardController) findJsapiService() {
	c.jsapiService = services_wxwork.NewJsapiService(c.app, c.wxOrg)
}

func (c *DashboardController) homeRedirectUri() string {
	u, _ := url.Parse(c.app.CallbackUrlHost())
	params := url.Values{}

	params.Add("corp_id", c.wxOrg.CorpId)
	params.Add("app_id", strconv.Itoa(c.app.Id))
	params.Add("agent_id", strconv.Itoa(c.app.AgentID))
	params.Add("mobile_redirect_uri", "")
	params.Add("pc_redirect_uri", "")

	u.Path = "/wxwork/home"
	u.RawQuery = params.Encode()

	return u.String()
}