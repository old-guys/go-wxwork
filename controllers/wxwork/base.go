package controllers_wxwork

import (
	"wxwork/controllers"
	"strings"
	"fmt"
)

type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	controllers.ApplicationController
}

func (c *BaseController) isMobile() (bv bool) {
	userAgent := strings.ToLower(c.UserAgent())
	devices := []string{"android", "iphone", "ios"}

	for _, device := range devices {
		if strings.Contains(userAgent, device) {
			bv = true
			break
		}
	}

	return bv
}

func (c *BaseController) isPc() bool {
	return ! c.isMobile()
}

func (c *BaseController) Prepare1() {
	fmt.Println(c.AppController.(string))
	//if app, ok := c.AppController.(NestPreparer); ok {
	//	app.NestPrepare()
	//}
}
