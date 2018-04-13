package models

import (
	"time"
	"wxwork/initializers"
	"wxwork/models/concerns"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

var DepartmentAr = initializers.DB.Model(&Department{})

type Department struct {
	Id        int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	OrgId     int       `gorm:"column:org_id;index;null"`
	Name      string    `gorm:"column:name;null"`
	ParentId  int       `gorm:"column:parent_id;index;null"`
	Status    int       `gorm:"column:status;null;default:1"`
	Path      string    `gorm:"column:path;index;null"`
	CreatedAt time.Time `gorm:"column:created_at;null"`
	UpdatedAt time.Time `gorm:"column:updated_at;null"`

	Parent interface{}  `gorm:"-"`

	UserDepartment        []UserDepartment
	UserAssistDepartments []UserAssistDepartment

	concerns.Tree
}

func (w Department) TableName() string {
	return "departments"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Department{})
}

func (w *Department) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}

func (w *Department) AfterSave(tx *gorm.DB) (err error) {
	paths := []string{"0", strconv.Itoa(w.Id)}

	if w.ParentId != 0 {
		if w.Parent == nil {
			department := Department{}
			DepartmentAr.Where(Department{Id: w.ParentId}).First(&department)
			w.Parent = department
		}
		paths = []string{w.Parent.(Department).Path, strconv.Itoa(w.Id)}
	}

	path := strings.Join(paths, "/")
	if w.Path != path {
		tx.Model(&w).UpdateColumn("path", strings.Join(paths, "/"))
	}

	return
}