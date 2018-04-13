package services_concerns_wxworkish

type MessageInterface struct {
	//Httpable
}

func (c *Base) SendCorpMsgWxwork(params map[string]interface{}) (data map[string]interface{}) {
	if c.App.AccessTokenExpired(c) { return data }

	return c.Post("cgi-bin/message/send?access_token=" + c.App.AccessToken, params)
}