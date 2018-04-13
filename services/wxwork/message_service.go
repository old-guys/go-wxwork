package services_wxwork

import (
	"github.com/imdario/mergo"
	"reflect"
	"strings"
	"wxwork/lib/wxwork"
	"wxwork/models"
	"wxwork/services/concerns/wxworkish"
)

type MessageService struct {
	services_concerns_wxworkish.Base
}

func NewMessageService(app models.WxworkApp) (MessageService) {
	service := &MessageService{}
	service.App = app
	//service.UpdateApp(app)
	return *service
}

func (c *MessageService) SendCorpMessage(opts map[string]interface{}) {
	headParams := c.GenerateMessageHead(opts)
	bodyParams := c.GenerateMessageBody(opts)

	mergo.Merge(&bodyParams, headParams)

	c.SendCorpMsgWxwork(bodyParams)
}

func (c *MessageService) GenerateMessageHead(opts map[string]interface{}) map[string]interface{} {
	userIds := opts["userIds"]
	deptIds := opts["deptIds"]
	tagIds := opts["tagIds"]

	if userIds!= nil && reflect.TypeOf(userIds).Kind() == reflect.Slice {
		userIds = strings.Join(userIds.([]string), "|")
	}

	if deptIds!= nil && reflect.TypeOf(deptIds).Kind() == reflect.Slice {
		deptIds = strings.Join(deptIds.([]string), "|")
	}

	if tagIds!= nil && reflect.TypeOf(tagIds).Kind() == reflect.Slice {
		tagIds = strings.Join(tagIds.([]string), "|")
	}

	return map[string]interface{}{
		"touser": userIds,
		"toparty": deptIds,
		"totag": tagIds,
		"agentid": opts["agent_id"],
	}
}

func (c *MessageService) GenerateMessageBody(opts map[string]interface{}) map[string]interface{} {
	msgType := opts["msg_type"]
	bodyParams := map[string]interface{}{}

	switch msgType {
		case "text":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"text": map[string]interface{}{
					"content": opts["content"],
				},
			})
		case "image":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"image": map[string]interface{}{
					"media_id": opts["media_id"],
				},
			})
		case "voice":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"voice": map[string]interface{}{
					"media_id": opts["media_id"],
				},
			})
		case "video":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"video": map[string]interface{}{
					"media_id": opts["media_id"],
					"title": opts["title"],
					"description": opts["description"],
				},
			})
		case "file":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"file": map[string]interface{}{
					"media_id": opts["media_id"],
				},
			})
		case "textcard":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"textcard": map[string]interface{}{
					"title": opts["title"],
					"description": opts["description"],
					"url": opts["url"],
				},
			})
		case "news":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"news": map[string]interface{}{
					"articles": opts["articles"],
				},
			})
		case "mpnews":
			mergo.Merge(&bodyParams, map[string]interface{}{
				"msgtype": msgType,
				"news": map[string]interface{}{
					"articles": opts["articles"],
				},
			})
		default:
			lib_wxwork.Logger.Info("Wx === generate_message_body: message type =", msgType, "can not support")
	}

	return bodyParams
}

/*
	service := services_wxwork.NewMessageService(app)
	service.SendCorpMessage(map[string]interface{}{
		"userIds": []string{"LiHui"},
		"msg_type": "textcard",
		"agent_id": "1000005",
		"title": "xxxx",
		"description": "xxxx",
		"url": "http://test.work.99zmall.com/wxwork/home",
	})
*/