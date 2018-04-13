package models

import (
	"wxwork/initializers"
	"time"
)

var WxworkUserTagMapAr = initializers.DB.Model(&WxworkUserTagMap{})

type WxworkUserTagMap struct {
	Id          int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID int       `gorm:"column:wxwork_org_id;index;null"`
	Userid      string    `gorm:"column:userid;index;null"`
	TagId       int       `gorm:"column:tag_id;index;null"`
	CreatedAt   time.Time `gorm:"column:created_at;null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;null"`
}

func (w WxworkUserTagMap) TableName() string {
	return "wxwork_user_tag_maps"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkUserTagMap{})
	initializers.DB.Model(&WxworkUserTagMap{}).AddUniqueIndex("uix_wxwork_org_id_userid_tag_id", "wxwork_org_id", "userid", "tag_id")
}

