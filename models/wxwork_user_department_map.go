package models

import (
	"time"
	"wxwork/initializers"
)

var WxworkUserDepartmentMapAr = initializers.DB.Model(&WxworkUserDepartmentMap{})

type WxworkUserDepartmentMap struct {
	Id                 int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID        int       `gorm:"column:wxwork_org_id;index;null"`
	WxworkUserId       int       `gorm:"column:wxwork_user_id;index;null"`
	WxworkDepartmentId int       `gorm:"column:wxwork_department_id;index;null"`
	CreatedAt          time.Time `gorm:"column:created_at;null"`
	UpdatedAt          time.Time `gorm:"column:updated_at;null"`
}

func (w WxworkUserDepartmentMap) TableName() string {
	return "wxwork_user_department_maps"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkUserDepartmentMap{})
}
