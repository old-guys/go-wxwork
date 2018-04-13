package controllers_api_wxwork

import (
	"wxwork/initializers"
	"wxwork/controllers"
)

type ApiWxBaseController struct {
	controllers.ApplicationController
}

func (c *ApiWxBaseController) Finish() {
	initializers.RemoveCaches()
}