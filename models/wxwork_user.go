package models

import (
	"time"
	"wxwork/initializers"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

var WxworkUserAr = initializers.DB.Model(&WxworkUser{})

type WxworkUser struct {
	Id             int    `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID    int    `gorm:"column:wxwork_org_id;index;null"`
	Userid         string `gorm:"column:userid;null"`
	Name           string `gorm:"column:name;null"`
	Mobile         string `gorm:"column:mobile;null"`
	Tel            string `gorm:"column:tel;null"`
	Email          string `gorm:"column:email;null"`
	Gender         int    `gorm:"column:gender;null"`
	Weixinid       string `gorm:"column:weixinid;null"`
	Position       string `gorm:"column:position;null"`
	Isleader       bool   `gorm:"column:isleader;null"`
	Avatar         string `gorm:"column:avatar;null"`
	EnglishName    string `gorm:"column:english_name;null"`
	DepartmentText string `gorm:"column:department;type:text;null"`
	UserType       int    `gorm:"column:user_type;null"`
	Status         int    `gorm:"column:status;default:1;null"`
	WxpluginStatus int    `gorm:"column:wxplugin_status;type:text;null"`
	Extattr        string `gorm:"column:extattr;type:text;null"`
	CreatedAt time.Time   `gorm:"column:created_at;null"`
	UpdatedAt time.Time   `gorm:"column:updated_at;null"`

	WxworkUserMap WxworkUserMap
	WxworkUserDepartmentMaps []WxworkUserDepartmentMap

	Department []interface{} `gorm:"-"`
}

func (w WxworkUser) TableName() string {
	return "wxwork_users"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkUser{})
	initializers.DB.Model(&WxworkUser{}).AddIndex("idx_wxwork_org_id_and_userid", "wxwork_org_id", "userid")
}

func (w *WxworkUser) BeforeSave(scope *gorm.Scope) (err error) {
	by, _ := yaml.Marshal(&w.Department)
	w.DepartmentText = string(by)
	scope.SetColumn("DepartmentText", string(by))

	return
}

func (w *WxworkUser) AfterFind() {
	yaml.Unmarshal([]byte(w.DepartmentText), &w.Department)
}

func assignAttributes(w interface{}, attrs map[string]interface{}) {
	for k, v := range attrs {
		initializers.DB.NewScope(w).SetColumn(k, v)
	}
}

func (w *WxworkUser) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}

func (w *WxworkUser) UserTicketExpired() bool {
	if w.Id == 0 { return true }

	if w.WxworkUserMap.Id == 0 {
		WxworkUserMapAr.Where(WxworkUserMap{WxworkUserId: w.Id}).First(&w.WxworkUserMap)
	}

	if w.WxworkUserMap.Id == 0 { return true }

	return w.WxworkUserMap.UserTicketExpired()
}

func (w *WxworkUser) User() (user User) {
	if w.Id == 0 { return user }

	if w.WxworkUserMap.Id == 0 {
		WxworkUserMapAr.Where(WxworkUserMap{WxworkUserId: w.Id}).First(&w.WxworkUserMap)
	}

	if w.WxworkUserMap.Id == 0 { return user }

	UserAr.Where(User{Id: w.WxworkUserMap.UserId}).First(&user)

	return user
}
