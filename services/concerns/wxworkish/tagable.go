package services_concerns_wxworkish

import (
	"wxwork/models"
	"strings"
	"strconv"
)

func (c *Base) TagAllUser(tagId int, userList *[]interface{}, more bool, refresh bool) {
	userData := c.GetUsersByTag(tagId)
	if userData["errcode"] == nil || userData["errcode"].(float64) != 0 {
		return
	}

	for _, v1 := range userData["userlist"].([]interface{}) {
		userData := v1.(map[string]interface{})
		data := map[string]interface{}{}

		for _, v2 := range *userList {
			info := v2.(map[string]interface{})
			if info["userid"] == userData["userid"] {
				data = info
			}
		}

		existing := data["userid"] == nil
		if !existing {
			data = c.GetUserInfo(userData["userid"].(string))
			if data["errcode"] == nil || data["errcode"].(float64) != 0 { continue }
		}

		if _, ok := data["tags"]; !ok { data["tags"] = []int{} }
		data["tags"] = append(data["tags"].([]int), tagId)

		if existing { continue }

		*userList = append(*userList, data)
	}
}

func (c *Base) TagsAllUsers(v interface{}, userList *[]interface{}, more bool, refresh bool) []interface{} {
	tagIds := []int{}

	switch v := v.(type) {
		case int:
			tagIds = []int{v}
		case []int:
			tagIds = v
	}

	if len(tagIds) == 0 {
		for _, id := range c.App.AllowTagIds {
			tagIds = append(tagIds, id.(int))
		}
	}

	if len(tagIds) == 0 { return  *userList }

	for _, id := range tagIds {
		c.TagAllUser(id, userList, more, refresh)
	}

	return *userList
}

func (c *Base) UpdateTag(data map[string]interface{}) {
	tagId, _ := strconv.Atoi(data["TagId"].(string))
	addUserItems := data["AddUserItems"].(string)
	delUserItems := data["DelUserItems"].(string)

	useridStr := ""
	if addUserItems != "" { useridStr = addUserItems }
	if delUserItems != "" { useridStr = delUserItems }

	userids := strings.Split(useridStr, ",")

	if addUserItems != "" { c.addUserItems(tagId, userids) }
	if delUserItems != "" { c.delUserItems(tagId, userids) }

	c.updateTagByWxUsers(userids)
}

func (c *Base) addUserItems(tagId int, userids []string) {
	if len(userids) == 0 { return }

	for _, userid := range userids {
		wutm := models.WxworkUserTagMap{WxworkOrgID: c.WxworkOrg.Id, Userid: userid, TagId: tagId}
		models.WxworkUserTagMapAr.Where(&wutm).FirstOrInit(&wutm)
		models.WxworkUserTagMapAr.Save(&wutm)
	}
}

func (c *Base) delUserItems(tagId int, userids []string) {
	if len(userids) == 0 { return }

	wutm := models.WxworkUserTagMap{WxworkOrgID: c.WxworkOrg.Id, TagId: tagId}
	models.WxworkUserTagMapAr.Where(&wutm).Delete(&wutm, "userid IN (?)", userids)
}

func (c *Base) updateTagByWxUsers(userids []string) {
	if len(userids) == 0 { return }

	for _, userid := range userids {
		c.UpdateUser(userid, map[string]interface{}{})
	}
}


