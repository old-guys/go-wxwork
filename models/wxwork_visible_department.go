package models

import (
	"wxwork/initializers"
)

var WxworkVisibleDepartmentAr = initializers.DB.Model(&WxworkVisibleDepartment{})

type WxworkVisibleDepartment struct {
	Id                 int `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID        int `gorm:"column:wxwork_org_id;index;null"`
	WxworkAppId        int `gorm:"column:wxwork_app_id;index;null"`
	WxworkDepartmentId int `gorm:"column:wxwork_department_id;index;null"`

	WxworkOrg        WxworkOrg
	WxworkApp        WxworkApp
	WxworkDepartment WxworkDepartment
}

func (w WxworkVisibleDepartment) TableName() string {
	return "wxwork_visible_departments"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkVisibleDepartment{})
	initializers.DB.Model(&WxworkVisibleDepartment{}).AddUniqueIndex("idx_org_id_app_id", "wxwork_org_id", "wxwork_app_id")
}
