package services_wxwork

import (
	"wxwork/services/concerns/wxworkish"
	"wxwork/models"
	"time"
	"strconv"
	"github.com/astaxie/beego/utils"
	"strings"
)

type JsapiService struct {
	services_concerns_wxworkish.Base
}

func NewJsapiService(app models.WxworkApp, wxOrg models.WxworkOrg) (service JsapiService) {
	service.App = app

	wxworkOrgMap := wxOrg.WxworkOrgMap
	org := models.Org{}
	if app.WxworkOrgID != wxOrg.Id {
		models.WxworkAppAr.Model(&app).Related(&wxOrg)
		models.WxworkOrgMapAr.Model(&wxOrg).Related(&wxworkOrgMap)
		models.OrgAr.Where(&models.Org{Id: wxworkOrgMap.OrgId}).Find(&org)
	}

	if wxworkOrgMap.Id == 0 && wxOrg.Id != 0 {
		models.WxworkOrgMapAr.Model(&wxOrg).Related(&wxworkOrgMap)
		models.OrgAr.Where(&models.Org{Id: wxworkOrgMap.OrgId}).Find(&org)
	}

	service.WxworkOrg = wxOrg
	service.Org = org

	return service
}

func (c *JsapiService) UserInfoByCode(code string) (wxUser models.WxworkUser) {
	data := c.GetUserInfoByCode(code)
	if data["errcode"] == nil || data["errcode"].(float64) != 0 || data["UserId"] == nil { return }

	return c.UserInfoByUserid(data["UserId"].(string), data)
}

func (c *JsapiService) UserInfoByUserid(userid string, data map[string]interface{}) (wxUser models.WxworkUser) {
	if userid == "" { return }

	models.WxworkUserAr.FirstOrInit(&wxUser, models.WxworkUser{Userid: userid})

	newWxUser := models.WxworkUser{}
	if (wxUser.Id != 0 && !wxUser.UserTicketExpired()) || (data["user_ticket"] != nil) {
		var expiredAt interface{}
		userTicket := wxUser.WxworkUserMap.UserTicket

		if data["user_ticket"] != nil {
			seconds := data["expires_in"].(float64)
			duration, _ := time.ParseDuration(strconv.FormatFloat(seconds, 'f', -1, 64) + "s")
			expiredAt = time.Now().Add(duration)

			userTicket = data["user_ticket"].(string)
		}

		newWxUser = c.UpdateUserDetail(userTicket, expiredAt, map[string]interface{}{})
	} else {
		newWxUser = c.UpdateUser(userid, map[string]interface{}{})
	}

	if newWxUser.Id != 0 { return newWxUser }

	return wxUser
}

func (c *JsapiService) GenerateConfig(fullpath string) (data map[string]interface{}) {
	if c.App.JsapiTicketExpired(c) { return data }

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := string(utils.RandomCreateBytes(8))
	requestURL := strings.Join([]string{c.App.CallbackUrlHost(), fullpath}, "")

	signature, _ := c.JsapiSign(map[string]interface{}{
		"noncestr": nonceStr,
		"timestamp": timestamp,
		"jsapi_ticket": c.App.JsapiTicket,
		"url": requestURL,
	})

	return map[string]interface{}{
		"appId": c.WxworkOrg.CorpId,
		"timestamp": timestamp,
		"nonceStr": nonceStr,
		"signature": signature,
		"expiredAt": c.App.JsapiTicketExpiredAt,
	}
}