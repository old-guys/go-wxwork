package controllers

import (
	"github.com/astaxie/beego"
	"strings"
)

type ApplicationController struct {
	beego.Controller
	ControllerName string
	ActionName string
}

type NestPreparer interface {
	NestPrepare()
}

func (c *ApplicationController) UserAgent() string {
	userAgent := c.Ctx.Request.Header["User-Agent"]

	return strings.Join(userAgent, "")
}

func (c *ApplicationController) MicroMessengerBrowser() bool {
	return strings.Contains(c.UserAgent(), "MicroMessenger")
}

func (c *ApplicationController) Prepare() {
	c.ControllerName, c.ActionName = c.GetControllerAndAction()

	c.Data["MicroMessengerBrowser"] = c.MicroMessengerBrowser()

	c.Layout = "index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["sharedWxwork"] = "shared/wxwork.html"

	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}
