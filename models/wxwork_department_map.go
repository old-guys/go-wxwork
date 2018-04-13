package models

import (
	"wxwork/initializers"
)

var WxworkDepartmentMapAr = initializers.DB.Model(&WxworkDepartmentMap{})

type WxworkDepartmentMap struct {
	Id                 int `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkDepartmentId int `gorm:"column:wxwork_department_id;index;null"`
	DepartmentId       int `gorm:"column:department_id;index;null"`

	WxworkDepartment   WxworkDepartment
	Department         Department
}

func (w WxworkDepartmentMap) TableName() string {
	return "wxwork_department_maps"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkDepartmentMap{})
}