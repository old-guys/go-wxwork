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
			result = c.eventListener(mapData)
		case "text":
		case "image":
		case "voice":
		case "video":
		case "location":
		case "link":
	}

	lib_wxwork.Logger.Info("Wxwork org === Send data to Wxwork:", result)
	return result
}

func (c *OrgService) eventListener(data map[string]interface{}) string {
	event := data["Event"]
	result := ""

	switch event {
		case "subscribe":
		case "unsubscribe":
		case "enter_agent":
		case "LOCATION":
		case "batch_job_result":
		case "change_contact":
			c.changeContact(data)
		case "click":
		case "view":
		case "scancode_push":
		case "scancode_waitmsg":
		case "pic_sysphoto":
		case "pic_photo_or_album":
		case "pic_weixin":
		case "location_select":
	}

	return result
}

func (c *OrgService) changeContact(data map[string]interface{})  {
	goworker.Enqueue(&goworker.Job{
		Queue: "myqueue",
		Payload: goworker.Payload{
			Class: "WxworkOrgChangeContact",
			Args: []interface{}{c.WxworkOrg.Id, data},
		},
	})
}