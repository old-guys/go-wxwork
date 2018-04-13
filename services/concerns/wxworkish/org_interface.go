package services_concerns_wxworkish

func (c *Base) GetOrgAccessToken() (data map[string]interface{}) {
	return c.Get("cgi-bin/gettoken?corpid=" + c.WxworkOrg.CorpId + "&corpsecret=" + c.WxworkOrg.Secret)
}
