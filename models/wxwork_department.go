package models

import (
	"time"
	"wxwork/initializers"
	"github.com/jinzhu/gorm"
	"strings"
	"strconv"
)

var WxworkDepartmentAr = initializers.DB.Model(&WxworkDepartment{})

type WxworkDepartment struct {
	Id           int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID  int       `gorm:"column:wxwork_org_id;index;null"`
	Name         string    `gorm:"column:name;null"`
	DeptId       int       `gorm:"column:dept_id;index;null"`
	DeptParentId int       `gorm:"column:dept_parent_id;index;null"`
	ParentId     int       `gorm:"column:parent_id;index;null"`
	Order        int       `gorm:"column:order;null"`
	Status       int       `gorm:"column:status;null;default:1"`
	Path         string    `gorm:"column:path;null"`
	CreatedAt    time.Time `gorm:"column:created_at;null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;null"`

	WxworkOrg               WxworkOrg
	//WxworkDepartmentMap     *WxworkDepartmentMap
	WxworkVisibleDepartment *WxworkVisibleDepartment
	WxworkUserDepartmentMaps []WxworkUserDepartmentMap

	Parent interface{} `gorm:"-"`
}

func (w WxworkDepartment) TableName() string {
	return "wxwork_departments"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkDepartment{})
	initializers.DB.Model(&WxworkDepartment{}).AddUniqueIndex("uix_wxwork_org_id_dept_id", "wxwork_org_id", "dept_id")
}

func (w WxworkDepartment) Department() (department Department) {
	wdMap := WxworkDepartmentMap{}
	WxworkDepartmentMapAr.Where("wxwork_department_id = ?", w.Id).First(&wdMap)

	if wdMap.Id != 0 {
		DepartmentAr.Where("id = ?", wdMap.DepartmentId).First(&department)
	}

	return department
}

func (w *WxworkDepartment) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}

func (w *WxworkDepartment) AfterSave(tx *gorm.DB) (err error) {
	paths := []string{"0", strconv.Itoa(w.Id)}

	if w.ParentId != 0 {
		if w.Parent == nil {
			WxworkDepartmentAr.Where(WxworkDepartment{Id: w.ParentId}).First(&w.Parent)
		}
		paths = []string{w.Parent.(WxworkDepartment).Path, strconv.Itoa(w.Id)}
	}

	path := strings.Join(paths, "/")
	if w.Path != path {
		tx.Model(&w).UpdateColumn("path", strings.Join(paths, "/"))
	}

	return
}