package news

import (
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsType struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	Name      string         `json:"name" validate:"required" gorm:"type:varchar(255)"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func NewsTypeList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	keys := make([]NewsType, 0)
	name := strings.ToLower(c.Query("name"))

	tx := db.Model(NewsType{}).Where("lower(name) LIKE ?", "%"+name+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&keys).Error
	if err != nil {
		return nil, err
	}
	p.Items = keys
	return p, nil
}

func (p *NewsType) Create() error {
	db := database.DBconn
	return db.Create(p).Error
}

func (p *NewsType) Update() error {
	db := database.DBconn
	return db.Updates(p).Error
}

func (p *NewsType) Delete() error {
	db := database.DBconn
	return db.Delete(p).Error
}

func GetName(id uint) string {
	db := database.DBconn
	var ct NewsType
	if err := db.First(&ct, id).Error; err != nil {
		return "Тодорхойгүй"
	}
	return ct.Name
}
