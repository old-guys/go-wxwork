package models

import (
	"time"
	"wxwork/initializers"
	"strconv"
	"gopkg.in/yaml.v2"
	"github.com/jinzhu/gorm"
	"net/url"
	"github.com/go-sql-driver/mysql"
)

var WxworkAppAr = initializers.DB.Model(&WxworkApp{})

type WxworkApp struct {
	Id                   int       `gorm:"column:id;AUTO_INCREMENT;primary_key"`
	WxworkOrgID          int       `gorm:"column:wxwork_org_id;index"`
	Name                 string    `gorm:"column:name;size(255);null"`
	AgentID              int       `gorm:"column:agent_id;default(0);null;"`
	Secret               string    `gorm:"column:secret;size(255);null"`
	Token                string    `gorm:"column:token;size(255);null"`
	AesKey               string    `gorm:"column:aes_key;size(255);null"`
	AccessToken          string	   `gorm:"column:access_token;size(1000);null"`
	//AccessTokenExpiredAt time.Time `gorm:"column:access_token_expired_at;type:datetime;null"`
	AccessTokenExpiredAt mysql.NullTime `gorm:"column:access_token_expired_at;type:datetime;null"`
	JsapiTicket          string	   `gorm:"column:jsapi_ticket;size(1000);null"`
	//JsapiTicketExpiredAt time.Time `gorm:"column:jsapi_ticket_expired_at;type:datetime;null"`
	JsapiTicketExpiredAt mysql.NullTime `sql:"column:jsapi_ticket_expired_at;type:datetime;null"`
	AllowUserinfosYaml   string    `gorm:"column:allow_userinfos;type:text;null"`
	AllowPartysYaml      string    `gorm:"column:allow_partys;type:text;null"`
	AllowTagsYaml        string    `gorm:"column:allow_tags;type:text;null"`
	LogoUrl              string    `gorm:"column:logo_url;size(255);null"`
	HomeUrl              string    `gorm:"column:home_url;size(255);null"`
	Description          string    `gorm:"column:description;size(255);null"`
	CallbackUrl          string    `gorm:"column:callback_url;size(255);null"`
	CreatedAt            time.Time `gorm:"column:created_at;null"`
	UpdatedAt            time.Time `gorm:"column:updated_at;null"`

	AllowUserinfos           map[string]interface{} `gorm:"-"`
	AllowPartys              map[string]interface{} `gorm:"-"`
	AllowTags                map[string]interface{} `gorm:"-"`
	AllowUserIds             []interface{}          `gorm:"-"`
	AllowPartyIds            []interface{}          `gorm:"-"`
	AllowTagIds              []interface{}          `gorm:"-"`

	WxworkOrg                WxworkOrg
	WxworkVisibleDepartments []WxworkVisibleDepartment
}

func (w WxworkApp) TableName() string {
	return "wxwork_apps"
}

type serviceInterface interface {
	GetAppAccessToken() (data map[string]interface{})
	GetJsapiTicket() (data map[string]interface{})
}

func init() {
	initializers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&WxworkApp{})
	initializers.DB.Model(&WxworkApp{}).AddUniqueIndex("uix_wxwork_org_id_agent_id", "wxwork_org_id", "agent_id")
}

func (w *WxworkApp) BeforeSave(scope *gorm.Scope) (err error) {
	by, _ := yaml.Marshal(&w.AllowUserinfos)
	w.AllowUserinfosYaml = string(by)
	scope.SetColumn("AllowUserinfosYaml", string(by))

	by, _ = yaml.Marshal(&w.AllowPartys)
	w.AllowPartysYaml = string(by)
	scope.SetColumn("AllowPartysYaml", string(by))

	by, _ = yaml.Marshal(&w.AllowTags)
	w.AllowTagsYaml = string(by)
	scope.SetColumn("AllowTagsYaml", string(by))

	return
}

func  (w *WxworkApp) AfterSave(scope *gorm.Scope) (err error) {
	w.assignAllow()

	return
}

func (w *WxworkApp) AfterFind() {
	w.assignAllow()
}

func (w *WxworkApp) assignAllow() {
	yaml.Unmarshal([]byte(w.AllowUserinfosYaml), &w.AllowUserinfos)
	yaml.Unmarshal([]byte(w.AllowPartysYaml), &w.AllowPartys)
	yaml.Unmarshal([]byte(w.AllowTagsYaml), &w.AllowTags)

	if w.AllowUserinfos["user"] != nil {
		for _, v := range w.AllowUserinfos["user"].([]interface{}) {
			w.AllowUserIds = append(w.AllowUserIds, v.(map[interface{}]interface{})["userid"])
		}
	}

	if w.AllowPartys["partyid"] != nil {
		w.AllowPartyIds = w.AllowPartys["partyid"].([]interface{})
	}

	if w.AllowTags["tagid"] != nil {
		w.AllowTagIds = w.AllowTags["tagid"].([]interface{})
	}
}

func (w *WxworkApp) AccessTokenExpired(service serviceInterface) bool {
	expired := time.Now().After(w.AccessTokenExpiredAt.Time)

	if w.AccessToken == "" || expired {
		w.UpdateAccessToken(service)
	}

	return time.Now().After(w.AccessTokenExpiredAt.Time)
}

func (w *WxworkApp) UpdateAccessToken(service serviceInterface) {

	data := service.GetAppAccessToken()

	if data["errcode"].(float64) == 0 {
		accessToken := data["access_token"].(string)

		seconds := data["expires_in"].(float64) - 300
		duration, _ := time.ParseDuration(strconv.FormatFloat(seconds, 'f', -1, 64) + "s")
		expiredAt := time.Now().Add(duration)

		initializers.DB.Model(&w).Update(WxworkApp{
			AccessToken: accessToken,
			AccessTokenExpiredAt: mysql.NullTime{ Time: expiredAt, Valid: true },
		})
	}
}

func (w *WxworkApp) AssignAttributes(attrs map[string]interface{}) {
	assignAttributes(w, attrs)
}

func (w *WxworkApp) JsapiTicketExpired(service serviceInterface) bool {
	expired := time.Now().After(w.JsapiTicketExpiredAt.Time)

	if w.JsapiTicket == "" || expired {
		w.UpdateJsapiTicket(service)
	}

	return time.Now().After(w.JsapiTicketExpiredAt.Time)
}

func (w *WxworkApp) UpdateJsapiTicket(service serviceInterface)  {
	data := service.GetJsapiTicket()

	if data["errcode"].(float64) == 0 {
		jsapiTicket := data["ticket"].(string)

		seconds := data["expires_in"].(float64) - 300
		duration, _ := time.ParseDuration(strconv.FormatFloat(seconds, 'f', -1, 64) + "s")
		expiredAt := time.Now().Add(duration)

		initializers.DB.Model(&w).Update(WxworkApp{
			JsapiTicket: jsapiTicket,
			JsapiTicketExpiredAt: mysql.NullTime{ Time: expiredAt, Valid: true },
		})
	}
}

func (w *WxworkApp) CallbackUrlHost() string {
	u, _ := url.Parse(w.CallbackUrl)
	u.Path = ""

	return u.String()
}