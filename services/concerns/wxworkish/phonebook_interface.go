package services_concerns_wxworkish

import (
	"encoding/json"
	"strconv"
	"fmt"
)

type PhonebookInterface struct {
	//Httpable
}

func (c *Base) GetCorpAgentInfo() (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	return c.Get("cgi-bin/agent/get?access_token=" + c.AccessToken() + "&agentid=" + fmt.Sprintf("%v", c.App.AgentID))
}

func (c *Base) GetDepartments(d interface{}) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	departmentId := ""
	if d != nil {
		by, _ := json.Marshal(&d)
		departmentId = string(by)
	}

	return c.Get("cgi-bin/department/list?access_token=" + c.AccessToken() + "&id=" + departmentId)
}

func (c *Base) GetUserInfo(userId string) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	return c.Get("cgi-bin/user/get?access_token=" + c.AccessToken() + "&userid=" + userId)
}

func (c *Base) GetSimpleUsersByDepartment(d int, f interface{}) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	departmentId := strconv.Itoa(d)

	fetchChild := "0"
	if f != nil {
		by, _ := json.Marshal(&f)
		fetchChild = string(by)
	}

	return c.Get("cgi-bin/user/simplelist?access_token=" + c.AccessToken() + "&department_id=" + departmentId + "&fetch_child=" + fetchChild)
}

func (c *Base) GetUsersByDepartment(d int, f interface{}) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	departmentId := strconv.Itoa(d)

	fetchChild := "0"
	if f != nil {
		by, _ := json.Marshal(&f)
		fetchChild = string(by)
	}

	return c.Get("cgi-bin/user/list?access_token=" + c.AccessToken() + "&department_id=" + departmentId + "&fetch_child=" + fetchChild)
}

func (c *Base) GetUsersByTag(tagId int) (data map[string]interface{}) {
	if c.AccessTokenExpired() { return data }

	return c.Get("cgi-bin/tag/get?access_token=" + c.AccessToken() + "&tagid=" + strconv.Itoa(tagId))
}
