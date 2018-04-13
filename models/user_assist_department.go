package models

import (
	"time"
	"wxwork/initializers"
)

var UserAssistDepartmentAr = initializers.DB.Model(&UserAssistDepartment{})

type UserAssistDepartment struct {
	Id           int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	UserId       int       `gorm:"column:user_id;index;null"`
	DepartmentId int       `gorm:"column:department_id;index;null"`
	CreatedAt    time.Time `gorm:"column:created_at;null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;null"`
}

func (w UserAssistDepartment) TableName() string {
	return "user_assist_departments"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&UserAssistDepartment{})
}
