package lang

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Lang struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;type:varchar(120)" validate:"required"`
	Name      string         `json:"name" gorm:"type:varchar(255)" validate:"required"`
	Image     string         `json:"image"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*Lang) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbl_languages"
}

func LangList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	langs := make([]Lang, 0)
	name := strings.ToLower(c.Query("name"))

	tx := db.Model(Lang{}).Where("lower(name) LIKE ?", "%"+name+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&langs).Error
	if err != nil {
		return nil, err
	}
	p.Items = langs
	return p, nil
}

func (p *Lang) Save() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Lang) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Lang) Remove() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}
