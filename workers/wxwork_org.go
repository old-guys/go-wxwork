package workers

import (
	"fmt"
	"wxwork/lib/wxwork"
	"wxwork/models"
	"wxwork/services/wxwork"
)

func changeContact(queue string, args ...interface{}) error {
	fmt.Printf("queue = %s, args = %v\n", queue, args)
	lib_wxwork.Logger.Info(queue, args)

	id := args[0]
	data := args[1].(map[string]interface{})
	wxOrg := models.WxworkOrg{}
	wxworkOrgMap := wxOrg.WxworkOrgMap
	org := models.Org{}

	models.WxworkOrgAr.Where("id = ?", id).First(&wxOrg)
	models.WxworkOrgMapAr.Model(&wxOrg).Related(&wxworkOrgMap)
	models.OrgAr.Where(&models.Org{Id: wxworkOrgMap.OrgId}).Find(&org)

	service := services_wxwork.OrgService{}
	service.WxworkOrg = wxOrg
	service.Org = org

	service.EventListener(data)


	return nil
}