package services_concerns_wxworkish

import (
	"wxwork/models"
	"sort"
	"strconv"
	"fmt"
	"wxwork/initializers"
	"github.com/imdario/mergo"
	"reflect"
	"regexp"
)

func (c *Base) DepartmentAllUsers(deptId int, refresh bool) (userlist []interface{}) {
	variableName := "@_all_user_infos_[" + strconv.Itoa(deptId) + "]"

	cache := initializers.GlobalCache.Get(variableName)
	if cache != nil && cache != "" {
		return cache.([]interface{})
	}

	data := c.GetUsersByDepartment(deptId, nil)
	if data["errcode"].(float64) != 0 { return userlist }

	userlist = data["userlist"].([]interface{})
	initializers.GlobalCache.Put(variableName, userlist, 0)

	return userlist
}

type BufferType string
const (
	Console BufferType = "console1"
	Channel BufferType = "channel"
	Conversation BufferType = "conversation"
)

func (c *Base) AllUserInfos(v interface{}, refresh bool) (userList []interface{}) {
	fmt.Println(BufferType(Console))

	deptIds := []float64{}

	switch v := v.(type) {
		case float64:
			deptIds = []float64{v}
		case float32:
			deptIds = []float64{float64(v)}
		case []float64:
			deptIds = v
	}

	userids := []string{}
	tagIds := []int{}

	if len(deptIds) == 0 {
		for _, k := range c.App.AllowPartyIds {
			deptIds = append(deptIds, float64(k.(int)))
		}

		for _, id := range c.App.AllowUserIds {
			userids = append(userids, id.(string))
		}

		for _, id := range c.App.AllowTagIds {
			tagIds = append(tagIds, id.(int))
		}
	}

	sort.Float64s(deptIds)
	key_prefix := "@_all_user_infos_:" + strconv.Itoa(c.App.Id) + ":"
	key := key_prefix + ""

	isContain, _ := Contain(float64(1), deptIds)
	if !isContain {
		key = key_prefix + fmt.Sprintf("%v", deptIds)
	}

	initializers.GlobalCache.Delete(key)
	cache := initializers.GlobalCache.Get(key)
	if cache != nil && cache != "" {
		return cache.([]interface{})
	}

	tags := map[string]interface{}{}
	deptIds = c.AllDepartmentDeptIds(deptIds, false)
	for _, k := range deptIds {
		data := c.DepartmentAllUsers(int(k), refresh)

		if len(data) == 0 { continue }
		for _, info := range data {
			str := info.(map[string]interface{})["userid"].(string)
			if tags[str] == true { continue }

			tags[str] = true
			userList = append(userList, info)
		}
	}

	for _, userid := range userids {
		if tags[userid] == true { continue }

		info := c.GetUserInfo(userid)
		if info["errcode"].(float64) != 0 { continue }

		userList = append(userList, info)
	}

	c.TagsAllUsers(tagIds, &userList, true, refresh)

	initializers.GlobalCache.Put(key, userList, 0)

	return userList
}

func (c * Base) AllUserIds(deptIds interface{}, refresh bool) (userids []string){
	tags := map[string]interface{}{}
	userList := c.AllUserInfos(deptIds, refresh)

	for _, info := range userList {
		userid := info.(map[string]interface{})["userid"].(string)
		if tags[userid] == true { continue }

		tags[userid] = true
		userids = append(userids, userid)
	}

	return userids
}

func (c *Base) UpdateAllUsers() {
	oldUserids := []string{}
	models.WxworkUserAr.Where("wxwork_org_id = ?", c.WxworkOrg.Id).Pluck("userid", &oldUserids)

	tags := map[string]interface{}{}
	allUserInfos := c.AllUserInfos(nil, false)

	for _, v := range allUserInfos {
		info := v.(map[string]interface{})
		info["skip_update_uz_department"] = true

		tags[info["userid"].(string)] = true

		c.UpdateUserFlow(info)
	}


	diffIds := []string{}
	for _, userid := range oldUserids {
		if tags[userid] != true { diffIds = append(diffIds, userid) }
	}

	c.HideUser(diffIds)
}

func (c *Base) UpdateDisabledUser(userids interface{}, status string) {

}

func (c *Base) HideUser(userids interface{}) {
	c.UpdateDisabledUser(userids, "hide")
}

func (c *Base) DeleteUser(userids interface{}) {
	c.UpdateDisabledUser(userids, "left")
}

func (c *Base) UpdateUser(userid string, extAttrs map[string]interface{}) (wxUser models.WxworkUser) {
	data := c.GetUserInfo(userid)
	if data["errcode"].(float64) != 0 { return wxUser }

	mergo.Merge(&data, extAttrs)
	return c.UpdateUserFlow(data)
}

func (c *Base) UpdateUserDetail(userTicket string, expiredAt interface{}, extAttrs map[string]interface{}) (wxUser models.WxworkUser)  {
	data := c.GetUserInfoDetail(userTicket)
	if data["errcode"].(float64) != 0 { return wxUser }

	mergo.Merge(&data, extAttrs)
	mergo.Merge(&data, map[string]interface{}{
		"user_ticket": userTicket,
		"user_ticket_expired_at": expiredAt,
	})

	return c.UpdateUserFlow(data)
}

func (c *Base) UpdateUserFlow(data map[string]interface{}) (wxUser models.WxworkUser) {
	userid, ok := data["userid"]
	if !ok || userid.(string) == "" { return wxUser }

	// 更新微信用户本身的属性
	wxUser, user := c.UpdateWxworkUserAttributes(data)

	// 更新微信用户与标签关系
	c.UpdateWxUserTagMaps(wxUser, data)

	// 更新微信用户与部门关系
	c.UpdateWxUserDepartmentMaps(&wxUser, nil)

	// 更新用户主、辅部门
	if value, ok := data["skip_update_uz_department"]; !ok || !value.(bool) {
		c.UpdateUserDepartmentMaps(wxUser, user, nil)
	}

	// 更新用户角色
	c.UpdateUserRole(wxUser, user, data)

	return wxUser
}

func (c *Base) UpdateWxworkUserAttributes(data map[string]interface{}) (wxUser models.WxworkUser, user models.User) {

	models.WxworkUserAr.FirstOrInit(&wxUser, map[string]interface{}{"userid": data["userid"], "wxwork_org_id": c.WxworkOrg.Id})
	attrs := c.AssignWxUserAttrs(data, wxUser)

	var err error
	tx := initializers.DB.Begin()

	if wxUser.Id != 0 {
		initializers.DB.Model(&wxUser).Related(&wxUser.WxworkUserMap)
		if wxUser.WxworkUserMap.Id != 0 {
			initializers.DB.Where(&models.User{Id: wxUser.WxworkUserMap.UserId}).First(&user)
		}
	}
	user.AssignAttributes(map[string]interface{}{"name": data["name"], "org_id": c.Org.Id})

	err = tx.Save(&user).Error
	if err != nil { tx.Callback(); return }

	wuMap := wxUser.WxworkUserMap
	wuMapAttrs := map[string]interface{}{"user_id": user.Id, "wxwork_org_id": c.WxworkOrg.Id, "id": wuMap.Id}
	mergo.Merge(&wuMapAttrs, attrs["WxworkUserMap"])
	wuMap.AssignAttributes(wuMapAttrs)

	attrs["WxworkUserMap"] = wuMap
	wxUser.AssignAttributes(attrs)

	err = tx.Save(&wxUser).Error
	if err != nil { tx.Callback(); return }

	tx.Commit()

	return wxUser, user
}

func (c *Base) UpdateWxUserTagMaps(wxUser models.WxworkUser, data map[string]interface{}) {
	if _, ok := data["tags"]; !ok { return }

	newTagIds := data["tags"].([]int)
	oldTagIds := []int{}
	models.WxworkUserTagMapAr.Where(&models.WxworkUserTagMap{Userid: wxUser.Userid}).Pluck("tag_id", &oldTagIds)

	tags1 := map[int]interface{}{}
	tags2 := map[int]interface{}{}
	for _, id := range newTagIds {
		tags1[id] = true
	}
	for _, id := range oldTagIds {
		tags2[id] = true
	}

	for _, id := range newTagIds {
		if tags2[id] == true { continue }
		models.WxworkUserTagMapAr.Create(&models.WxworkUserTagMap{WxworkOrgID: c.Org.Id, Userid: wxUser.Userid, TagId: id})
	}

	diff := []int{}
	for _, id := range oldTagIds {
		if tags1[id] == true { continue }
		diff = append(diff, id)
	}

	if len(diff) > 0 {
		wutm := models.WxworkUserTagMap{WxworkOrgID: c.Org.Id, Userid: wxUser.Userid}
		models.WxworkUserTagMapAr.Where(&wutm).Delete(&wutm, "tag_id IN (?)", diff)
	}
}

func (c *Base) UpdateWxUserDepartmentMaps(wxUser *models.WxworkUser, deptIds interface{}) {
	if deptIds == nil { deptIds = wxUser.Department }

	wudm := models.WxworkUserDepartmentMap{WxworkUserId: wxUser.Id, WxworkOrgID: wxUser.WxworkOrgID}
	mapDepartmentIds := []int{}
	models.WxworkUserDepartmentMapAr.Where(&wudm).Pluck("wxwork_department_id", &mapDepartmentIds)

	wxDepartmentIds := []int{}
	models.WxworkDepartmentAr.Where("dept_id IN (?) AND wxwork_org_id = ?", deptIds, c.WxworkOrg.Id).Pluck("id", &wxDepartmentIds)

	tags1 := map[int]interface{}{}
	tags2 := map[int]interface{}{}
	for _, id := range mapDepartmentIds {
		tags1[id] = true
	}
	for _, id := range wxDepartmentIds {
		tags2[id] = true
	}

	for _, id := range wxDepartmentIds {
		if tags1[id] == true { continue }
		wudmap := models.WxworkUserDepartmentMap{WxworkOrgID: wxUser.WxworkOrgID, WxworkUserId: wxUser.Id, WxworkDepartmentId: id}
		models.WxworkUserDepartmentMapAr.Create(&wudmap)
	}

	diff := []int{}
	for _, id := range mapDepartmentIds {
		if tags2[id] == true { continue }
		diff = append(diff, id)
	}
	if len(diff) > 0 {
		models.WxworkUserDepartmentMapAr.Where(&wudm).Delete(&wudm, "wxwork_department_id IN (?)", diff)
	}
}

func (c *Base) UpdateUserDepartmentMaps(wxUser models.WxworkUser, user models.User, deptIds interface{}) {
	if deptIds == nil { deptIds = wxUser.Department }

	visibleUserids := c.App.AllowUserIds
	visibleDeptIds := c.App.AllowPartyIds
	visibleTagIds := c.App.AllowTagIds

	if len(visibleUserids) > 0 && len(visibleDeptIds) == 0 {
		isContain, _ := Contain(visibleUserids, wxUser.Userid)
		if !isContain {
			//Contains
		}
	}

	type Department struct {
		Id int
		DeptId int
	}
	var departments []Department
	models.DepartmentAr.Table("departments").
		Joins("INNER JOIN `wxwork_department_maps` ON `wxwork_department_maps`.`department_id` = `departments`.`id` INNER JOIN `wxwork_departments` ON `wxwork_departments`.`id` = `wxwork_department_maps`.`wxwork_department_id`").
		Select("`departments`.`id`, `wxwork_departments`.`dept_id` as dept_id").
		Where("`departments`.`org_id` = ? AND `wxwork_departments`.`dept_id` IN (?) AND `wxwork_departments`.`status` IN (?)", c.Org.Id, deptIds, []int{0, 1}).
		Find(&departments)

	if departments != nil {
		// 主部门
		department := Department{}
		for _, v := range deptIds.([]interface{}) {
			var deptId int

			switch v := v.(type) {
				case float64:
					deptId = int(v)
				case float32:
					deptIds = int(v)
				case int:
					deptIds = v
			}

			if department.Id != 0 { break }

			for _, d := range departments {
				if d.DeptId == deptId {
					department = d
					break
				}
			}

		}
		user.Status = 1
		models.UserAr.Model(&user).Related(&user.UserDepartment)
		user.UserDepartment = models.UserDepartment{Id: user.UserDepartment.Id, DepartmentId: department.Id}

		// 辅部门
		usersAssistDepartmentIds := []int{}
		models.UserAssistDepartmentAr.Where(models.UserAssistDepartment{UserId: user.Id}).Pluck("department_id", &usersAssistDepartmentIds)
		assistDepartmentIds := []int{}
		for _, d := range departments {
			if d.Id == department.Id { continue }
			assistDepartmentIds = append(assistDepartmentIds, d.Id)
		}

		tags1 := map[int]interface{}{}
		tags2 := map[int]interface{}{}
		for _, id := range usersAssistDepartmentIds {
			tags1[id] = true
		}
		for _, id := range assistDepartmentIds {
			tags2[id] = true
		}

		for _, id := range assistDepartmentIds {
			if tags1[id] == true { continue }

			uad := models.UserAssistDepartment{UserId: user.Id, DepartmentId: id}
			models.UserAssistDepartmentAr.Create(&uad)
		}

		diff := []int{}
		for  _, id := range usersAssistDepartmentIds {
			if tags2[id] == true { continue }
			diff = append(diff, id)
		}
		if len(diff) > 0 {
			uad := models.UserAssistDepartment{UserId: user.Id}
			models.UserAssistDepartmentAr.Where(&uad).Delete(&uad, "department_id IN (?)", diff)
		}

	} else {
		ud := models.UserDepartment{UserId: user.Id}
		models.UserDepartmentAr.Where(&ud).Delete(&ud)

		uad := models.UserAssistDepartment{UserId: user.Id}
		models.UserAssistDepartmentAr.Where(&uad).Delete(&uad)

		user.Status = 1
		is, _ := Contain(wxUser.Userid, visibleUserids)
		if !is { user.Status = 0 }

		if !is && len(visibleTagIds) > 0 {
			wutm := models.WxworkUserTagMap{Userid: wxUser.Userid}
			models.WxworkUserTagMapAr.Where(wutm).Where("tag_id IN (?)", visibleTagIds).First(&wutm)
			if wutm.Id != 0 { user.Status = 1 }
		}
	}

	models.UserAr.Save(&user)
}

func (c *Base) UpdateUserRole(wxUser models.WxworkUser, user models.User, data map[string]interface{}) {

}

func (c *Base) AssignWxUserAttrs(data map[string]interface{}, wxUser models.WxworkUser) (attrs map[string]interface{}) {
	attrs = map[string]interface{}{ "wxwork_org_id": c.WxworkOrg.Id }

	assignOriginKeys := map[string]interface{}{
		"userid": "userid",
		"name": "name",
		"mobile": "mobile",
		"tel": "tel",
		"weixinid": "weixinid",
		"email": "email",
		"gender": "gender",
		"position": "position",
		"isleader": "isleader",
		"avatar": "avatar",
		"english_name": "english_name",
		"extattr": "extattr",
		"wxplugin_status": "wxplugin_status",
		"status": "status",
		"department": "department",
		"user_type": "usertype",
	}

	lambdas := map[string]interface{}{
		"gender": func(gender string) int {
			i, _ := strconv.Atoi(gender)
			return i
		},
		"avatar": func(avatar string) string {
			return regexp.MustCompile(`^http:\/\/`).ReplaceAllString(avatar, "https://")
		},
		"isleader": func(isleader float64) int {
			return 1
		},
		"user_type": func(userType int) int {
			return 1
		},
		"wxplugin_status": func(wxpluginStatus int) int {
			return 1
		},
	}

	for assignKey, originKey := range assignOriginKeys {
		if _, ok := data[originKey.(string)]; !ok { continue }
		attrs[assignKey] = data[originKey.(string)]

		if _, ok := lambdas[assignKey]; !ok { continue }

		fv := reflect.ValueOf(lambdas[assignKey])
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(attrs[assignKey])
		rs := fv.Call(params)
		attrs[assignKey] = rs[0].Interface()
	}

	if data["user_ticket"] != nil && data["user_ticket_expired_at"] != nil {
		if wxUser.Id != 0 { initializers.DB.Model(&wxUser).Related(&wxUser.WxworkUserMap) }
		wuMap := wxUser.WxworkUserMap

		if wuMap.UserTicket != data["user_ticket"].(string) {
			attrs["WxworkUserMap"] = map[string]interface{}{
				"wxwork_org_id": c.WxworkOrg.Id,
				"user_ticket": data["user_ticket"],
				"user_ticket_expired_at": data["user_ticket_expired_at"],
			}
		}
	}

	fmt.Println("attrs", attrs)

	return attrs
}