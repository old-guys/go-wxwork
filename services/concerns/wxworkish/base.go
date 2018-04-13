package services_concerns_wxworkish

import (
	"wxwork/models"
	//"fmt"
)

type Base struct {
	App models.WxworkApp
	WxworkOrg models.WxworkOrg
	Org models.Org
}

func (c *Base) AccessTokenExpired() (bool) {
	if c.WxworkOrg.EnabledBookSyn { return c.WxworkOrg.AccessTokenExpired(c) }

	return c.App.AccessTokenExpired(c)
}

func (c *Base) AccessToken() string {
	if c.WxworkOrg.EnabledBookSyn { return  c.WxworkOrg.AccessToken }

	return c.App.AccessToken
}