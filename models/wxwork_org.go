package models

import (
	"time"
	"wxwork/initializers"
	"github.com/go-sql-driver/mysql"
	"strconv"
)

var WxworkOrgAr = initializers.DB.Model(&WxworkOrg{})

type WxworkOrg struct {
	Id                   int            `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	CorpName             string         `gorm:"column:corp_name;size(255);null"`
	CorpId               string         `gorm:"column:corp_id;index;unique_index;null"`
	Secret               string         `gorm:"column:secret;size(255);null"`
	Token                string         `gorm:"column:token;size(255);null"`
	AesKey               string         `gorm:"column:aes_key;size(255);null"`
	AccessToken          string	        `gorm:"column:access_token;size(1000);null"`
	AccessTokenExpiredAt mysql.NullTime `gorm:"column:access_token_expired_at;type:datetime;null"`
	EnabledBookSyn       bool           `gorm:"column:enabled_book_syn;null"`
	CreatedAt            time.Time      `gorm:"column:created_at;null"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;null"`

	WxworkApps               []WxworkApp
	WxworkDepartment         []WxworkDepartment
	WxworkVisibleDepartments []WxworkVisibleDepartment
	WxworkOrgMap WxworkOrgMap
}

type OrgService interface {
	GetOrgAccessToken() (data map[string]interface{})
}

func (w WxworkOrg) TableName() string {
	return "wxwork_orgs"
}

func init()  {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkOrg{})
}

func (w *WxworkOrg) AccessTokenExpired(service OrgService) bool {
	expired := time.Now().After(w.AccessTokenExpiredAt.Time)

	if w.AccessToken == "" || expired {
		w.UpdateAccessToken(service)
	}

	return time.Now().After(w.AccessTokenExpiredAt.Time)
}

func (w *WxworkOrg) UpdateAccessToken(service OrgService) {

	data := service.GetOrgAccessToken()

	if data["errcode"].(float64) == 0 {
		accessToken := data["access_token"].(string)

		seconds := data["expires_in"].(float64) - 300
		duration, _ := time.ParseDuration(strconv.FormatFloat(seconds, 'f', -1, 64) + "s")
		expiredAt := time.Now().Add(duration)
		initializers.DB.Model(&w).Update(WxworkOrg{
			AccessToken: accessToken,
			AccessTokenExpiredAt: mysql.NullTime{ Time: expiredAt, Valid: true },
		})
	}
}

