package controllers_api_wxwork

import (
	"wxwork/services/wxwork"
	"wxwork/models"
	"wxwork/lib/wxwork"
	"regexp"
	"encoding/xml"
)

type ApiWxOrgController struct {
	ApiWxBaseController
	service services_wxwork.OrgService
}

type msgXml struct {
	ToUserName string `xml:"ToUserName"`
	AgentID string `xml:"AgentID"`
	Encrypt string `xml:"Encrypt"`
}

func (c *ApiWxOrgController) Prepare() {
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		id = c.Ctx.Input.Query("id")
	}

	if id != ""  {
		models.WxworkOrgAr.Where("id = ?", id).First(&c.service.WxworkOrg)
	}
}

func (c *ApiWxOrgController) Callback() {
	result := "failure"
	if c.service.WxworkOrg.Id == 0 {
		lib_wxwork.Logger.Info("wxwork === WxworkOrg not exist, params =")
		c.Ctx.Output.Body([]byte(result))
	}

	encrypt := c.GetString("echostr")
	mx := msgXml{}

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
		"aes_key": c.service.WxworkOrg.AesKey,
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

func (c *ApiWxOrgController) validate(encrypt string) bool {
	nonce := c.GetString("nonce")
	timestamp := c.GetString("timestamp")

	msgSignature, err := c.service.Sign(map[string]interface{}{
		"token": c.service.WxworkOrg.Token,
		"nonce": nonce,
		"timestamp": timestamp,
		"encrypt": encrypt,
	})

	if err != nil { return false }

	return msgSignature == c.GetString("msg_signature")
}

