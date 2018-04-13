package services_wxwork

import (
	"wxwork/services/concerns/wxworkish"
	"wxwork/models"
	"wxwork/lib/wxwork"
	"reflect"
	"wxwork/initializers"
	"encoding/json"
	"github.com/benmanns/goworker"
)

type OrgService struct {
	services_concerns_wxworkish.Base
}

func NewOrgService(app models.WxworkApp) (service OrgService) {
	service.App = app

	wxworkOrg := app.WxworkOrg
	wxworkOrgMap := wxworkOrg.WxworkOrgMap
	org := models.Org{}
	if app.WxworkOrgID != wxworkOrg.Id {
		models.WxworkAppAr.Model(&app).Related(&wxworkOrg)
		models.WxworkOrgMapAr.Model(&wxworkOrg).Related(&wxworkOrgMap)
		models.OrgAr.Where(&models.Org{Id: wxworkOrgMap.OrgId}).Find(&org)
	}
	service.WxworkOrg = wxworkOrg
	service.Org = org

	return service
}

func (c *OrgService) MsgTypeListener(data interface{}) string {
	lib_wxwork.Logger.Info("Wxwork org === msg_type_listener:", data)

	if reflect.TypeOf(data).Kind() == reflect.String {
		return data.(string)
	}

	redis := initializers.Redis
	result := "success"
	mapData := data.(map[string]interface{})

	value, _ := json.Marshal(mapData)
	key := string(value)

	ok, _ := redis.Do("GET", key)
	if ok != nil { return result }

	redis.Do("SET", key, 1, "EX", 5)
	msgType := mapData["MsgType"]

	switch msgType {
		case "event":
			c.addChangeContactQueue(mapData)
	}

	lib_wxwork.Logger.Info("Wxwork org === Send data to Wxwork:", result)
	return result
}

func (c *OrgService) EventListener(data map[string]interface{}) string {
	event := data["Event"]
	result := "success"

	switch event {
		case "change_contact":
			c.changeContact(data)
	}

	return result
}

func (c *OrgService) addChangeContactQueue(data map[string]interface{}) {
	goworker.Enqueue(&goworker.Job{
		Queue: "myqueue",
		Payload: goworker.Payload{
			Class: "WxworkOrgChangeContact",
			Args: []interface{}{c.WxworkOrg.Id, data},
		},
	})
}

func (c *OrgService) changeContact(data map[string]interface{}) {
	changeType := data["ChangeType"]

	switch changeType {
		case "create_user", "update_user":
			c.updateUsersHandle(data)
		case "delete_user":
			c.deleteUsersHandle(data)
		case "create_party", "update_party":
			c.updateDeptsHandle(data)
		case "delete_party":
			c.deleteDeptsHandle(data)
		case "update_tag":
			c.updateTagsHandle(data)
	}

}

func (c *OrgService) updateUsersHandle(data map[string]interface{}) {
	userid := data["UserID"].(string)
	newUserid, ok := data["NewUserID"].(string)

	if ok {
		models.WxworkUserAr.
			Where(models.WxworkUser{WxworkOrgID: c.WxworkOrg.Id}).
			Where("userid = ?", userid).
			Updates(models.WxworkUser{Userid: newUserid})
	}

	c.UpdateUser(userid, map[string]interface{}{})
}

func (c *OrgService) deleteUsersHandle(data map[string]interface{}) {
	userid := data["UserID"]
	c.DeleteUser(userid)
}

func (c *OrgService) updateDeptsHandle(data map[string]interface{}) {
	deptId := data["Id"]
	c.UpdateDepartment(deptId)
}

func (c *OrgService) deleteDeptsHandle(data map[string]interface{}) {
	deptId := data["Id"]
	c.DeleteDepartment(deptId)
}

func (c *OrgService) updateTagsHandle(data map[string]interface{}) {
	c.UpdateTag(data)
}
