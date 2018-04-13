package services_concerns_wxworkish

type AppInterface struct {
	//Httpable
}

func (c *Base) GetAppAccessToken() (data map[string]interface{}) {
	return c.Get("cgi-bin/gettoken?corpid=" + c.WxworkOrg.CorpId + "&corpsecret=" + c.App.Secret)
}