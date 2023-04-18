package organization

import (
	"encoding/json"
	"os"
	"strings"

	"crowdfund/pkg/base/role"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/convertor"

	"gorm.io/datatypes"
)

type RoleResponse struct {
	Id      uint                     `json:"id,string"`
	Name    string                   `json:"name"`
	Modules []map[string]interface{} `json:"modules"`
	Actions []UserAction             `json:"actions"`
	Role    *UserRole                `json:"role,omitempty"`
}

type UserRole struct {
	Id      uint           `json:"id,string"`
	OrgId   uint           `json:"org_id,string"`
	Name    string         `json:"name"`
	Actions datatypes.JSON `json:"actions"`
}

type UserModule struct {
	Id          uint   `json:"id,string"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
}

type UserPage struct {
	Id        uint   `json:"id,string"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Sequence  uint   `json:"sequence"`
	GroupCode string `json:"group_code"`
	GroupName string `json:"group_name"`
}
type UserAction struct {
	Id          uint   `json:"id,string"`
	ModuleId    uint   `json:"module_id,string"`
	ModuleName  string `json:"module_name"`
	PageId      uint   `json:"page_id,string"`
	PageName    string `json:"page_name" `
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *RoleResponse) Init(orgId, userId uint) {
	db := database.DBconn
	err := db.Table(os.Getenv("DB_SCHEMA")+".tbr_roles as r").Select("r.id, ou.org_id, r.name, r.actions").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_organization_users as ou ON (r.id = ou.role_id AND r.deleted_at is NULL)").Where("ou.org_id=? AND ou.user_id=?", orgId, userId).Take(&r.Role).Error

	if err == nil {
		str := string(r.Role.Actions)

		str = strings.Replace(str, "[", "", -1)
		str = strings.Replace(str, "]", "", -1)
		str = strings.ReplaceAll(str, " ", "")
		str = strings.ReplaceAll(str, `"`, "")
		str = strings.ReplaceAll(str, "\"", "")

		strArr := strings.Split(str, ",")
		var actionIds []int
		for _, val := range strArr {
			numVal := convertor.StringToInt(val)
			actionIds = append(actionIds, numVal)
		}

		if len(actionIds) != 0 {
			// db.Model(role.Action{}).Find(&r.Actions, actionIds)

			tx := db.Table(os.Getenv("DB_SCHEMA")+".tbr_actions as a").Select("a.*").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbr_pages as p ON (a.page_id = p.id AND p.deleted_at is NULL)").Where("a.id IN ?", actionIds)

			tx.Find(&r.Actions)
		}

		var moduleIds []int
		var pageIds []int

		for _, action := range r.Actions {
			moduleIds = append(moduleIds, int(action.ModuleId))
			pageIds = append(pageIds, int(action.PageId))
		}

		moduleIds = helpers.UniqueIntSlice(moduleIds)
		pageIds = helpers.UniqueIntSlice(pageIds)
		var modules []UserModule
		db.Model(role.Module{}).Order("sequence").Find(&modules, moduleIds)

		m, _ := json.Marshal(&modules)
		_ = json.Unmarshal(m, &r.Modules)

		for _, mod := range r.Modules {
			var pages []UserPage
			db.Model(role.Page{}).Order("sequence").Where("module_id = ? AND id IN ?", mod["id"], pageIds).Find(&pages)
			mod["pages"] = pages
		}
	} else {
		r.Actions = []UserAction{}
		r.Modules = make([]map[string]interface{}, 0)
	}
}

func (r *RoleResponse) Return() {
	r.Id = r.Role.Id
	r.Name = r.Role.Name
	r.Role = nil
}

func GetUserRoleWithActions(orgId, userId uint) *RoleResponse {
	roleRes := new(RoleResponse)
	roleRes.Init(orgId, userId)
	roleRes.Return()
	return roleRes
}
