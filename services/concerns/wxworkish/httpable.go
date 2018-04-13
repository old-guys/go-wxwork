package services_concerns_wxworkish

import (
	"github.com/astaxie/beego"
	"wxwork/lib/wxwork"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"time"
	"crypto/tls"
)

var wxworkQyapiHost = beego.AppConfig.String("wxwork_qyapi_host")

func (c *Base) Get(path string) (data map[string]interface{}) {
	url := strings.Join([]string{wxworkQyapiHost, path}, "/")

	start_time := time.Now().UnixNano()
	req := httplib.Get(url)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.Header("Content-Type", "application/json")

	body, err := req.String()
	exec_time := time.Now().UnixNano() - start_time

	json.Unmarshal([]byte(body), &data)

	lib_wxwork.Logger.Info("get_data: url:", url, "exec_time:", exec_time, "body =", body, err)

	return data
}

func (c *Base) Post(path string, params map[string]interface{}) (data map[string]interface{}) {

	url := strings.Join([]string{wxworkQyapiHost, path}, "/")
	jsonBytes, err := json.Marshal(params)

	start_time := time.Now().UnixNano()
	req := httplib.Post(url)
	req.JSONBody(params)
	req.Header("Content-Type", "application/json")

	//res, err := req.Response()
	//body, err := ioutil.ReadAll(res.Body)
	body, err := req.String()
	exec_time := time.Now().UnixNano() - start_time
	json.Unmarshal([]byte(body), &data)

	lib_wxwork.Logger.Info("post_json: url:", url, "exec_time:", exec_time, "params =", string(jsonBytes), "body =", body, err)

	return data
}