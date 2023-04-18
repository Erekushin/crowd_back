package role

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Role struct {
	Id          uint           `json:"id,string" gorm:"primaryKey"`
	OrgId       uint           `json:"org_id,string" validate:"required"`
	Name        string         `json:"name" gorm:"type:varchar(120)" validate:"required"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	Actions     datatypes.JSON `json:"actions" validate:"required"`
	CreatedAt   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"created_by,omitempty"`
	UpdatedAt   data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy   uint           `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy   uint           `json:"-"`
}

func (*Role) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbr_roles"
}

func RoleList(c *fiber.Ctx, orgId int) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	roles := make([]Role, 0)
	name := strings.ToLower(c.Query("name"))

	tx := db.Model(Role{}).Where("org_id=?", orgId).Where("lower(name) LIKE ?", "%"+name+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	p.Items = roles
	return p, nil
}

func (r *Role) Create() error {
	db := database.DBconn

	if err := db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Role) Update() error {
	db := database.DBconn
	if err := db.Updates(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Role) Delete() error {
	db := database.DBconn
	if err := db.Delete(&r).Error; err != nil {
		return err
	}

	return nil
}
