package services_concerns_wxworkish

import (
	"wxwork/lib/wxwork"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"time"
	"strconv"
)

type Cipherable struct {
	Key string
	AesKey string
	Token string
}

/*
	services_concerns_wxworkish.Base{}.Decrypt(map[string]interface{}{
		"aes_key": "xxxx",
		"key": "xxxx",
		"data": "xxxx",
	})
*/
func (c *Base) Decrypt(args map[string]interface{}) (interface{}, error) {
	aes_key := args["aes_key"]
	if aes_key == nil {
		aes_key = beego.AppConfig.String("wxwork_aes_key")
	}

	key := args["key"]
	if key == nil {
		key = beego.AppConfig.String("wxwork_key")
	}

	params := map[string]interface{}{
		"aes_key": aes_key,
		"key": key,
		"data": args["data"],
	}

	return lib_wxwork.Cipher.Decrypt(params)
}

/*
	services_concerns_wxworkish.Base{}.Encrypt(map[string]interface{}{
		"aes_key": "xxxx",
		"key": "xxxx",
		"params": "xxxx",
		"params": map[string]interface{}{
			"name": "xxxx"
		},
	})
*/
func (c *Base) Encrypt(args map[string]interface{}) (string, error) {
	aes_key := args["aes_key"]
	if aes_key == nil {
		aes_key = beego.AppConfig.String("wxwork_aes_key")
	}

	key := args["key"]
	if key == nil {
		key = beego.AppConfig.String("wxwork_key")
	}

	params := map[string]interface{}{
		"aes_key": aes_key,
		"key": key,
		"params": args["params"],
	}

	return lib_wxwork.Cipher.Encrypt(params)
}

/*
	services_concerns_wxworkish.Base{}.Sign(map[string]interface{}{
		"nonce": "xxxx",
		"timestamp": "xxxx",
		"token": "xxxx",
		"encrypt": "xxxx",
	})
*/
func (c *Base) Sign(args map[string]interface{}) (string, error) {
	nonce := args["nonce"]
	if nonce == nil {
		nonce = string(utils.RandomCreateBytes(8))
	}

	timestamp := args["timestamp"]
	if timestamp == nil {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	token := args["token"]
	if token == nil {
		token = beego.AppConfig.String("wxwork_token")
	}

	opts := map[string]interface{}{
		"nonce": nonce,
		"timestamp": timestamp,
		"token": token,
		"encrypt": args["encrypt"],
	}

	return lib_wxwork.Cipher.Sign(opts)
}

/*
	services_concerns_wxworkish.Base{}.JsapiSign(map[string]interface{}{
		"noncestr": "xxxx",
		"timestamp": "xxxx",
		"jsapi_ticket": "xxxx",
		"url": "xxxx",
	})
*/
func (c *Base) JsapiSign(args map[string]interface{}) (string, error) {
	noncestr := args["noncestr"]
	if noncestr == nil {
		noncestr = string(utils.RandomCreateBytes(8))
	}

	timestamp := args["timestamp"]
	if timestamp == nil {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	opts := map[string]interface{}{
		"noncestr": noncestr,
		"timestamp": timestamp,
		"jsapi_ticket": args["jsapi_ticket"],
		"url": args["url"],
	}

	return lib_wxwork.Cipher.JsapiSign(opts)
}