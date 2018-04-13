package services_wxwork

import (
	"github.com/astaxie/beego/utils"
	"wxwork/services/concerns/wxworkish"
	"time"
	"strconv"
	"reflect"
	"wxwork/initializers"
	"encoding/json"
	"wxwork/lib/wxwork"
)

type AppService struct {
	services_concerns_wxworkish.Base
}

func (c *AppService) MsgTypeListener(data interface{}) string {
	lib_wxwork.Logger.Info("Wxwork app === msg_type_listener:", data)

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
			result = c.EventListener(mapData)
		case "text":
		case "image":
		case "voice":
		case "video":
		case "location":
		case "link":
	}

	lib_wxwork.Logger.Info("Wxwork app === Send data to Wxwork:", result)
	return result
}

func (c *AppService) EventListener(data map[string]interface{}) string {
	event := data["Event"]
	result := ""

	switch event {
		case "subscribe":
		case "unsubscribe":
		case "enter_agent":
			result = c.enterAgent()
		case "LOCATION":
		case "batch_job_result":
		case "change_contact":
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

func (c *AppService) enterAgent() string {
	str := `<xml><ToUserName><![CDATA[LiHui]]></ToUserName><FromUserName><![CDATA[wwea672916e0a3a7c4]]></FromUserName><CreateTime>` + strconv.FormatInt(time.Now().Unix(), 10) + `</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[欢迎进入应用!]]></Content></xml>`

	msgEncrypt, err := c.Encrypt(map[string]interface{}{
		"params": str,
	})
	if err != nil { return "" }

	nonce := string(utils.RandomCreateBytes(8))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	msgSignature, err := c.Sign(map[string]interface{}{
		"nonce": nonce,
		"timestamp": timestamp,
		"encrypt": msgEncrypt,
	})
	if err != nil { return "" }

	return `<xml><Encrypt><![CDATA[` + msgEncrypt + `]]></Encrypt><MsgSignature><![CDATA[` + msgSignature + `]]></MsgSignature><TimeStamp>` + timestamp + `</TimeStamp><Nonce><![CDATA[` + nonce + `]]></Nonce></xml>`
}