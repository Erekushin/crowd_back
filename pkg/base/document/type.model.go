package document

import (
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Type struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"unique;not null"`
	Name      string         `json:"name" gorm:"type:varchar(100)"`
	NameEn    string         `json:"name_en" gorm:"type:varchar(100)"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*Type) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_document_types"
}

func TypeList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	keys := make([]Type, 0)

	tx := db.Model(Type{})

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&keys).Error
	if err != nil {
		return nil, err
	}
	p.Items = keys
	return p, nil
}

func (p *Type) Create() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Type) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Type) Delete() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}
