package services_concerns_wxworkish


type JsapiInterface struct {
	//Httpable
}

func (c *Base) GetJsapiTicket() (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	return c.Get("cgi-bin/get_jsapi_ticket?access_token=" + c.AccessToken())
}