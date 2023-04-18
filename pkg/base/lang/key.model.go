package lang

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Key struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;type:varchar(120)" validate:"required"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*Key) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbl_keys"
}

func KeyList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	keys := make([]Key, 0)
	code := strings.ToLower(c.Query("code"))

	tx := db.Model(Key{}).Where("lower(code) LIKE ?", "%"+code+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Order("code").Offset(p.Offset).Limit(p.PageSize).Find(&keys).Error
	if err != nil {
		return nil, err
	}
	p.Items = keys
	return p, nil
}

func (p *Key) Save() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Key) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Key) Remove() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}

	if err := db.Unscoped().Delete(&Translation{}, "key_id=?", p.Id).Error; err != nil {
		return err
	}
	return nil
}
