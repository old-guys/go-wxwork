package controllers_api_wxwork

import (
	"encoding/xml"
	"wxwork/services/wxwork"
	"wxwork/initializers"
	"wxwork/lib/wxwork"
	"regexp"
)

type ApiWxAppsController struct {
	ApiWxBaseController
	service services_wxwork.AppService
}

type MsgXml struct {
	ToUserName string `xml:"ToUserName"`
	AgentID string `xml:"AgentID"`
	Encrypt string `xml:"Encrypt"`
}

func (c *ApiWxAppsController) Prepare() {
	appId := c.Ctx.Input.Param(":id")
	if appId == "" {
		appId = c.Ctx.Input.Query("id")
	}

	if appId != ""  {
		initializers.DB.Where("id = ?", appId).Find(&c.service.App)//.Related(&c.service.App.WxworkOrg)
	}
}

func (c *ApiWxAppsController) Callback() {
	result := "failure"
	lib_wxwork.Logger.Info("c.service.App", c.service.App)
	if c.service.App.Id == 0 {
		lib_wxwork.Logger.Info("wxwork === App not exist, params =")
		c.Ctx.Output.Body([]byte(result))
	}

	encrypt := c.GetString("echostr")
	mx := MsgXml{}

	if len(encrypt) == 0 {
		encrypt = string(c.Ctx.Input.RequestBody)
	}

	reg := regexp.MustCompile(`^<xml+`)
	isXml := reg.MatchString(encrypt)
	if isXml {
		err := xml.Unmarshal([]byte(encrypt), &mx)
		if err == nil { encrypt = mx.Encrypt }
	}

	data, err := c.service.Decrypt(map[string]interface{}{
		"aes_key": c.service.App.AesKey,
		"key": c.service.WxworkOrg.CorpId,
		"data": encrypt,
	})
	lib_wxwork.Logger.Info("wxwork === Decrypted data:", data, "err =", err)

	if err != nil || !c.validate(encrypt) {
		lib_wxwork.Logger.Info("wxwork === Cipher failure: data =", data, ", params = #{params}")
		c.Ctx.Output.Body([]byte(result))
	}

	result = c.service.MsgTypeListener(data)

	c.Ctx.Output.Body([]byte(result))
}

func (c *ApiWxAppsController) SynOrg() {
	service := services_wxwork.NewOrgService(c.service.App)
	service.ActiveApp()

	c.Ctx.Output.Body([]byte("success"))
}

func (c *ApiWxAppsController) validate(encrypt string) bool {
	nonce := c.GetString("nonce")
	timestamp := c.GetString("timestamp")

	msgSignature, err := c.service.Sign(map[string]interface{}{
		"token": c.service.App.Token,
		"nonce": nonce,
		"timestamp": timestamp,
		"encrypt": encrypt,
	})

	if err != nil { return false }

	return msgSignature == c.GetString("msg_signature")
}

