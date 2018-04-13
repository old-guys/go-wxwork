package models

import (
	"time"
	"wxwork/initializers"
	// "wxwork/models/concerns"
)

var UserAr = initializers.DB.Model(&User{})

type User struct {
	Id        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	OrgId     int       `gorm:"column:org_id;index;null"`
	Name      string    `gorm:"column:name;null"`
	Status    int       `gorm:"column:status;default:1;null"`
	CreatedAt time.Time `gorm:"column:created_at;null"`
	UpdatedAt time.Time `gorm:"column:updated_at;null"`

	UserDepartment        UserDepartment
	UserAssistDepartments []UserAssistDepartment

	//concerns.Tree
}

func (w User) TableName() string {
	return "users"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
}

func (w *User) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}