package services_concerns_wxworkish

import (
	"wxwork/lib/wxwork"
	"strconv"
	"wxwork/initializers"
	"fmt"
	"reflect"
	"sort"
	"errors"
	"wxwork/models"
	"github.com/imdario/mergo"
	"strings"
)

const DEFAULT_MAX_DEPTH = 6

func (c *Base) MaxDeepLength() int {
	return DEFAULT_MAX_DEPTH
}

func (c *Base) UpdateDepartment(depatrmentId interface{}) {
	data := c.GetDepartments(depatrmentId)

	if data["errcode"].(float64) != 0 { return }

	departmentInfo := map[string]interface{}{}
	for _, v := range data["department"].([]interface{}) {
		dept := v.(map[string]interface{})

		if dept["id"].(float64) == depatrmentId.(float64) {
			departmentInfo = dept
			break
		}
	}

	c.UpdateDepartmentInternal(departmentInfo)
}

func (c *Base) UpdateDisabledDepartments(v interface{}, status string) {
	
}

func (c *Base) DeleteDepartment(deptIds interface{}) {
	c.UpdateDisabledDepartments(deptIds, "left")
}

func (c *Base) HideDepartment(deptIds interface{}) {
	c.UpdateDisabledDepartments(deptIds, "hide")
}

func (c *Base) AllDepartmentDeptIds(deptIds interface{}, refresh bool) (ids []float64) {
	departmentInfos := c.WxDepartmentInfos(deptIds, refresh)
	if departmentInfos == nil { return ids }

	for _, v := range departmentInfos {
		ids = append(ids, v.(map[string]interface{})["id"].(float64))
	}

	return ids
}

func (c *Base) WxDepartmentInfos(v interface{}, refresh bool) (data []interface{}) {
	deptIds := []float64{}

	switch v := v.(type) {
		case float64:
			deptIds = []float64{v}
		case float32:
			deptIds = []float64{float64(v)}
		case []float64:
			deptIds = v
	}

	if len(deptIds) == 0 {
		partyIds := c.App.AllowPartyIds
		for _, k := range partyIds {
			deptIds = append(deptIds, float64(k.(int)))
		}

		allowUserinfos := c.App.AllowUserIds
		if len(allowUserinfos) == 0 && len(partyIds) == 0 {
			deptIds = []float64{1}
		}
	}

	if len(deptIds) == 0 { return data }

	sort.Float64s(deptIds)

	key_prefix := "wxwork_dept_ids:" + strconv.Itoa(c.App.Id) + ":"
	key := key_prefix + "1"

	isContain, _ := Contain(float64(1), deptIds)
	if !isContain {
		key = key_prefix + fmt.Sprintf("%v", deptIds)
	}

	cache := initializers.GlobalCache.Get(key)
	if cache != nil && cache != "" {
		return cache.([]interface{})
	}

	tags := map[string]interface{}{}
	for _, k := range deptIds {
		if tags["1"] == true { break }

		str := strconv.FormatFloat(k, 'f', -1, 64)
		if tags[str] == true { continue }

		tags[str] = true
		departmentInfo := c.GetDepartments(k)
		if departmentInfo["errcode"].(float64) != 0 { continue }

		deptKey := key_prefix + str
		initializers.GlobalCache.Put(deptKey, departmentInfo["department"], 0)

		for _, info := range departmentInfo["department"].([]interface{}) {
			data = append(data, info)
		}
	}

	initializers.GlobalCache.Put(key, data, 0)

	return data
}

func (c *Base) RegroupAllDepartments(sortDepts []interface{}) []interface{} {
	var rootDept interface{}

	for _, v := range sortDepts {
		if v.(map[string]interface{})["parentid"].(float64) == 0 {
			rootDept = v
		}
	}

	if rootDept == nil {
		rootDept = map[string]interface{}{
			"id": float64(1),
			"name": "xxx",
			"parentid": float64(0),
			"order": 0,
			"new_parentid": float64(0),
		}

		sortDepts = append(sortDepts, rootDept)
	} else {
		for _, v := range sortDepts {
			info := v.(map[string]interface{})

			if info["parentid"] == nil || info["parentid"].(float64) == 0 {
				info["new_parentid"] = float64(0)
			}

			v = info
		}
	}

	for _, v1 := range sortDepts {
		dept1 := v1.(map[string]interface{})
		if _, ok := dept1["new_parentid"]; ok { continue }

		var parentDept interface{}
		for _, v2 := range sortDepts {
			dept2 := v2.(map[string]interface{})
			if dept1["parentid"] == dept2["id"] {
				parentDept = v2
			}
		}

		if parentDept == nil {
			dept1["new_parentid"] = rootDept.(map[string]interface{})["parentid"]
		} else {
			dept1["new_parentid"] = dept1["parentid"]
		}

		v1 = dept1
	}

	return c.SortAllDepartments(sortDepts, []interface{}{}, &[]interface{}{})
}

func (c *Base) SortAllDepartments(sortDepts []interface{}, parentDepts []interface{}, depts *[]interface{}) (data []interface{}) {
	if len(parentDepts) == 0 {
		for _, v := range sortDepts {
			dept := v.(map[string]interface{})

			if dept["new_parentid"].(float64) == 0 {
				parentDepts = append(parentDepts, dept)
			}
		}
	}

	if len(parentDepts) == 0 { return data }
	for _, v := range parentDepts { *depts = append(*depts, v) }

	childDepts := []interface{}{}
	for _, v1 := range parentDepts {
		parentDept := v1.(map[string]interface{})

		for _, v2 := range sortDepts {
			dept := v2.(map[string]interface{})

			if dept["new_parentid"].(float64) == parentDept["id"].(float64) {
				childDepts = append(childDepts, dept)
			}
		}

		if len(childDepts) != 0 {
			c.SortAllDepartments(sortDepts, childDepts, depts)
		}
	}

	return *depts
}

func (c *Base) UpdateAllDepartmens() {
	departmentInfos := c.WxDepartmentInfos(nil, false)
	if len(departmentInfos) == 0 { return }

	lib_wxwork.Logger.Info("update_all_departments department_infos =", departmentInfos)

	departmentInfos = c.RegroupAllDepartments(departmentInfos)
	lib_wxwork.Logger.Info("update_all_departments regroup_all_departments department_infos =", departmentInfos)
	if len(departmentInfos) == 0 { return }

	oldDeptIds := []float64{}
	models.WxworkDepartmentAr.Where("wxwork_org_id = ?", c.WxworkOrg.Id).Pluck("dept_id", &oldDeptIds)

	tags := map[string]interface{}{}

	for _, v := range departmentInfos {
		info := v.(map[string]interface{})
		info["skip_update_visible_department"] = true
		info["update_all"] = true

		tags[strconv.FormatFloat(info["id"].(float64), 'f', -1, 64)] = true

		c.UpdateDepartmentInternal(info)
	}

	diffIds := []float64{}
	for _, v := range oldDeptIds {
		str := strconv.FormatFloat(v, 'f', -1, 64)
		if tags[str] != true { diffIds = append(diffIds, v) }
	}

	c.HideDepartment(diffIds)
}

func (c *Base) UpdateDepartmentInternal(data map[string]interface{}) {
	if data["id"] == nil { return }

	attrs := map[string]interface{}{
		"name": data["name"],
		"dept_id": data["id"],
		"dept_parent_id": data["parentid"],
		"order": data["order"],
		"wxwork_org_id": c.WxworkOrg.Id,
	}

	if data["update_all"] == nil && data["parentid"].(float64) != 0 {
		var count int
		models.WxworkDepartmentAr.
			Where("wxwork_org_id = ? AND dept_id = ?", c.WxworkOrg.Id, data["parentid"]).
			Count(&count)

		if count == 0 { c.UpdateDepartment(data["parentid"]) }
	}

	wxParentDepartment := models.WxworkDepartment{}
	if _, ok := data["new_parentid"]; ok {
		models.WxworkDepartmentAr.Where("wxwork_org_id = ? AND dept_id = ?", c.WxworkOrg.Id, data["new_parentid"]).First(&wxParentDepartment)
	} else {

		if _, ok := data["parentid"]; ok {
			models.WxworkDepartmentAr.Where("wxwork_org_id = ? AND dept_id = ?", c.WxworkOrg.Id, data["parentid"]).First(&wxParentDepartment)
		} else {
			models.WxworkDepartmentAr.Where("wxwork_org_id = ? AND parent_id IS NULL", c.WxworkOrg.Id).First(&wxParentDepartment)
		}
	}

	wxDepartment := models.WxworkDepartment{}
	models.WxworkDepartmentAr.Where("wxwork_org_id = ? AND dept_id = ?", c.WxworkOrg.Id, data["id"]).FirstOrInit(&wxDepartment, attrs)
	isPersisted := wxDepartment.Id == 0
	deptParentIdChanged := !(wxDepartment.DeptParentId == wxParentDepartment.DeptId)

	department := wxDepartment.Department()
	parentDepartment := c.FindParentByDepth(wxParentDepartment.Department(),nil)

	lib_wxwork.Logger.Info("update_department_internal wx_department:", wxDepartment)
	lib_wxwork.Logger.Info("update_department_internal wx_parent_department:", wxParentDepartment)

	tx := initializers.DB.Begin()
	var err error
	department.AssignAttributes(map[string]interface{}{"name": data["name"], "parent_id": parentDepartment.Id, "org_id": c.Org.Id})
	err = tx.Save(&department).Error
	//if department.Id == 0 {
	//	department = models.Department{Name: data["name"].(string), ParentId: parentDepartment.Id}
	//	err = tx.Create(&department).Error
	//} else {
	//	err = tx.Model(&department).Update(map[string]interface{}{"name": data["name"], "parent_id": parentDepartment.Id}).Error
	//}

	if err != nil {
		tx.Rollback()
		return
	}

	lib_wxwork.Logger.Info("update_department_internal department:", department)
	lib_wxwork.Logger.Info("update_department_internal parent_department:", parentDepartment)

	lib_wxwork.Logger.Info("update_department_internal wxDepartment:", wxDepartment)
	lib_wxwork.Logger.Info("update_department_internal wxParentDepartment:", wxParentDepartment)

	mergo.Merge(&attrs, map[string]interface{}{"parent_id": wxParentDepartment.Id, "status": 1, "Parent": wxParentDepartment})
	wxDepartment.AssignAttributes(attrs)
	err = tx.Save(&wxDepartment).Error
	//if wxDepartment.Id == 0 {
	//	err = tx.Create(&wxDepartment).Error
	//} else {
	//	err = tx.Model(&wxDepartment).Update(attrs).Error
	//}

	if err != nil {
		tx.Rollback()
		return
	}

	wdMap := models.WxworkDepartmentMap{WxworkDepartmentId: wxDepartment.Id, DepartmentId: department.Id}
	err = initializers.DB.FirstOrInit(&wdMap, wdMap).Save(&wdMap).Error

	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	if data["skip_update_visible_department"] != true {
		c.UpdateVisibleDepartment(wxParentDepartment, wxDepartment, isPersisted, deptParentIdChanged)
	}
}

func (c *Base) UpdateVisibleDepartment(wxParentDepartment models.WxworkDepartment, wxDepartment models.WxworkDepartment, isPersisted bool, deptParentIdChanged bool) {
	if !deptParentIdChanged { return }

}

func (c *Base) FindParentByDepth(department models.Department, depth interface{}) models.Department {
	if department.Id == 0 { return department }
	if depth == nil { depth = c.MaxDeepLength() }

	if depth.(int) < 1 { return department }

	paths := strings.Split(department.Path, "/")
	if len(paths) - 1 < depth.(int) { return department }

	models.DepartmentAr.Where(&models.Department{OrgId: c.Org.Id}).Where(map[string]interface{}{"id": paths[1]}).First(&department)

	return department
}

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

