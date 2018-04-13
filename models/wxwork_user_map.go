package models

import (
	"wxwork/initializers"
	"time"
)

var WxworkUserMapAr = initializers.DB.Model(&WxworkUserMap{})

type WxworkUserMap struct {
	Id                  int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID         int       `gorm:"column:wxwork_org_id;index;null"`
	WxworkUserId        int       `gorm:"column:wxwork_user_id;index;null"`
	UserId              int       `gorm:"column:user_id;index;null"`
	UserTicket          string    `gorm:"column:user_ticket;size(1000);null"`
	UserTicketExpiredAt time.Time `gorm:"column:user_ticket_expired_at;null"`
}

func (w WxworkUserMap) TableName() string {
	return "wxwork_user_maps"
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkUserMap{})
}

func (w *WxworkUserMap) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}

func (w *WxworkUserMap) UserTicketExpired() bool {
	if len(w.UserTicket) == 0 { return true }

	return time.Now().After(w.UserTicketExpiredAt)
}