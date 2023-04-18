package organization

import (
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrgType struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(200)" validate:"required"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*OrgType) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organization_types"
}

func GetTypeName(id uint) string {
	db := database.DBconn
	var ct OrgType
	if err := db.First(&ct, id).Error; err != nil {
		return "Тодорхойгүй"
	}
	return ct.Name
}

func OrgTypeList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	modules := make([]OrgType, 0)

	tx := db.Model(OrgType{})
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&modules).Error
	if err != nil {
		return nil, err
	}
	p.Items = modules
	return p, nil
}
