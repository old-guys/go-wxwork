package models

import (
	"wxwork/initializers"
	"time"
)

var WxworkOrgMapAr = initializers.DB.Model(&WxworkOrgMap{})

type WxworkOrgMap struct {
	Id          int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgId int       `gorm:"column:wxwork_org_id;index;null"`
	OrgId       int       `gorm:"column:org_id;index;null"`
	CreatedAt   time.Time `gorm:"column:created_at;null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;null"`
}

func (w WxworkOrgMap) TableName() string {
	return "wxwork_org_maps"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkOrgMap{})
}