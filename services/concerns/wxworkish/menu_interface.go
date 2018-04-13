package services_concerns_wxworkish

type MenuInterface struct {
	//Httpable
}

func (c *Base) CreateMenu(params map[string]interface{}) (data map[string]interface{}) {
	if c.App.AccessTokenExpired(c) { return data }

	return c.Post("cgi-bin/menu/create?access_token=" + c.App.AccessToken + "&agentid=" + string(c.App.AgentID), params)
}

func (c *Base) GetMenu() (data map[string]interface{}) {
	if c.App.AccessTokenExpired(c) { return data }

	return c.Get("cgi-bin/menu/get?access_token=" + c.App.AccessToken + "&agentid=" + string(c.App.AgentID))
}

func (c *Base) DeleteMenu() (data map[string]interface{}) {
	if c.App.AccessTokenExpired(c) { return data }

	return c.Get("cgi-bin/menu/delete?access_token=" + c.App.AccessToken + "&agentid=" + string(c.App.AgentID))
}

