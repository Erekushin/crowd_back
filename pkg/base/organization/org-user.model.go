package organization

import (
	"os"

	"crowdfund/pkg/base/role"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
)

type OrgUser struct {
	OrgId     uint           `json:"org_id,string" validate:"required"`
	UserId    uint           `json:"user_id,string" validate:"required"`
	RoleId    int            `json:"role_id,string" gorm:"default:-1"`
	RoleName  string         `json:"role_name" gorm:"type:varchar(120)"`
	CreatedAt data.LocalTime `gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
}

func (*OrgUser) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organization_users"
}

type UsersResponse struct {
	Id           uint           `json:"id,string"`
	RegNo        string         `json:"reg_no"`
	FamilyName   string         `json:"family_name"`
	LastName     string         `json:"last_name"`
	FirstName    string         `json:"first_name"`
	Username     string         `json:"username"`
	Gender       int            `json:"gender"`
	BirthDate    string         `json:"birth_date"`
	Email        string         `json:"email"`
	PhoneNo      string         `json:"phone_no"`
	IsForeign    uint           `json:"is_foreign"`
	ProfileImage string         `json:"profile_image"`
	CountryName  string         `json:"country_name,omitempty"`
	RoleId       int            `json:"role_id,string"`
	RoleName     string         `json:"role_name"`
	CreatedAt    data.LocalTime `json:"created_date"`
	CreatedBy    uint           `json:"created_by"`
}

func UserList(c *fiber.Ctx, orgId int) (*data.Pagination, error) {
	var (
		err      error
		totalRow int64
	)
	result := make([]UsersResponse, 0)
	db := database.DBconn
	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_users as u").Select("u.*,ou.role_id, ou.role_name, ou.created_at, ou.created_by").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_organization_users as ou ON (u.id = ou.user_id AND u.deleted_at is NULL)").Where("ou.org_id=?", orgId)
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err = tx.Offset(p.Offset).Limit(p.PageSize).Find(&result).Error
	if err != nil {
		return nil, err
	}
	p.Items = result
	return p, nil
}

func (r *OrgUser) Create() error {
	db := database.DBconn
	role := role.Role{}
	db.First(&role, "id=?", r.RoleId)
	r.RoleName = role.Name
	return db.Create(&r).Error
}

func (r *OrgUser) SetRole() error {
	db := database.DBconn

	if r.RoleId > 0 {
		role := role.Role{}
		db.First(&role, "id=?", r.RoleId)
		r.RoleName = role.Name
	} else {
		r.RoleName = ""
	}

	// return db.Model(OrgUser{}).Where("org_id=? AND user_id=?", r.OrgId, r.UserId).Update("role_id", r.RoleId).Error
	return db.Model(OrgUser{}).Where("org_id=? AND user_id=?", r.OrgId, r.UserId).Updates(r).Error
}

func (r *OrgUser) Add() error {
	db := database.DBconn
	var cnt int64

	db.Model(OrgUser{}).Where("org_id=? AND user_id=?", r.OrgId, r.UserId).Count(&cnt)

	if cnt > 0 {
		return r.SetRole()
	}

	return r.Create()

}

func (r *OrgUser) Remove() error {
	db := database.DBconn
	return db.Where("org_id=?", r.OrgId).Where("user_id=?", r.UserId).Delete(&r).Error
}
