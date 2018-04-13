package services_concerns_wxworkish

type LoginInterface struct {
	//Httpable
}

func (c *Base) GetUserInfoByCode(code string) (data map[string]interface{}) {
	if c.App.AccessTokenExpired(c) { return data }

	return c.Get("cgi-bin/user/getuserinfo?access_token=" + c.App.AccessToken + "&code=" + code)
}

func (c *Base) GetUserInfoDetail(userTicket string) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	params := map[string]interface{}{
		"user_ticket": userTicket,
	}

	return c.Post("cgi-bin/user/getuserdetail?access_token=" + c.AccessToken(), params)
}