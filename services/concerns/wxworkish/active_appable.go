package services_concerns_wxworkish

import (
	"wxwork/lib/wxwork"
	"time"
	"wxwork/initializers"
)

func (c *Base) ActiveApp() {
	// 更新自建应用信息
	startTime := time.Now().UnixNano()
	c.UpdateAppInfo()
	endTime := time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_app(更新自建应用信息) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)

	// 同步组织架构信息
	startTime = time.Now().UnixNano()
	c.UpdateOrg()
	endTime = time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_org(同步组织架构) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)

	// 同步业务可见范围
	startTime = time.Now().UnixNano()
	c.UpdateVisibleScopes()
	endTime = time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_visible_scopes(同步可见范围) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)
}

// 更新自建应用信息
func (c *Base) UpdateAppInfo(){
	data := c.GetCorpAgentInfo()

	//c.App.AllowUserinfos = data["allow_userinfos"].(map[string]interface{})
	//c.App.AllowPartys = data["allow_partys"].(map[string]interface{})
	//c.App.AllowTags = data["allow_tags"].(map[string]interface{})
	c.App.AssignAttributes(map[string]interface{}{
		"AllowUserinfos": data["allow_userinfos"],
		"AllowPartys": data["allow_partys"],
		"AllowTags": data["allow_tags"],
	})
	initializers.DB.Model(&c.App).Update("id", c.App.Id)
}

// 同步组织架构
func (c *Base) UpdateOrg() {
	lib_wxwork.Logger.Info("update_org: corp_id =", c.WxworkOrg.CorpId)

	// 同步微信组织架构中的所有部门
	startTime := time.Now().UnixNano()
	c.UpdateAllDepartmens()
	endTime := time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_all_departments(同步企业微信所有部门) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)

	// 更新企业微信可见部门列表
	startTime = time.Now().UnixNano()
	c.UpdateWxVisibleDepartments()
	endTime = time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_wx_visible_departments(更新企业微信可见部门列表) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)

	// 同步微信组织架构中的所有用户
	startTime = time.Now().UnixNano()
	c.UpdateAllUsers()
	endTime = time.Now().UnixNano()
	lib_wxwork.Logger.Info("active_app update_all_users(同步微信组织架构中的所有用户) start_time:", startTime, "end_time:", endTime, "exec_time:", endTime - startTime)
}

// 更新企业微信可见部门列表
func (c *Base) UpdateWxVisibleDepartments() {
	lib_wxwork.Logger.Info("update_organization: ")
}

// 同步业务可见范围
func (c *Base) UpdateVisibleScopes() {

}