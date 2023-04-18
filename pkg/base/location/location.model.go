package location

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Location struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	OrgId     uint           `json:"org_id,string" validate:"required"`
	Name      string         `json:"name" gorm:"type:varchar(500)" validate:"required"`
	AimagCode string         `json:"aimag_code" gorm:"type:varchar(5)" validate:"required"`
	AimagName string         `json:"aimag_name" gorm:"type:varchar(255)"`
	SumCode   string         `json:"sum_code" gorm:"type:varchar(5)" validate:"required"`
	SumName   string         `json:"sum_name" gorm:"type:varchar(255)"`
	BagCode   string         `json:"bag_code" gorm:"type:varchar(5)" validate:"required"`
	BagName   string         `json:"bag_name" gorm:"type:varchar(255)"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*Location) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_locations"
}

func LocationList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	locations := make([]Location, 0)
	name := strings.ToLower(c.Query("name"))

	tx := db.Model(Location{}).Where("lower(name) LIKE ?", "%"+name+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&locations).Error
	if err != nil {
		return nil, err
	}
	p.Items = locations
	return p, nil
}

func (p *Location) Create() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Location) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Location) Remove() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}
