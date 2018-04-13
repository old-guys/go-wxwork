package models

import (
	"time"
	"wxwork/initializers"
)

var OrgAr = initializers.DB.Model(&Org{})

type Org struct {
	Id        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	Name      string    `gorm:"column:name;null"`
	CreatedAt time.Time `gorm:"column:created_at;null"`
	UpdatedAt time.Time `gorm:"column:updated_at;null"`

}

func (w Org) TableName() string {
	return "orgs"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Org{})
}