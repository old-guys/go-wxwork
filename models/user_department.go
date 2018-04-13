package models

import (
	"time"
	"wxwork/initializers"
)

var UserDepartmentAr = initializers.DB.Model(&UserDepartment{})

type UserDepartment struct {
	Id           int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	UserId       int       `gorm:"column:user_id;index;null"`
	DepartmentId int       `gorm:"column:department_id;index;null"`
	CreatedAt    time.Time `gorm:"column:created_at;null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;null"`
}

func (w UserDepartment) TableName() string {
	return "user_departments"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&UserDepartment{})
}