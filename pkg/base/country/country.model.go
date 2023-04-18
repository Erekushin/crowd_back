package country

import (
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Country struct {
	Code          uint           `json:"code"`
	IsoAlphaCode3 string         `json:"iso_alpha_code_3" gorm:"column:iso_alpha_code_3"`
	IsoAlphaCode2 string         `json:"iso_alpha_code_2" gorm:"column:iso_alpha_code_2"`
	CommonName    string         `json:"common_name"`
	FullName      string         `json:"full_name"`
	EnName        string         `json:"en_name"`
	Status        string         `json:"status"`
	CreatedAt     data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy     uint           `json:"created_by,omitempty"`
	UpdatedAt     data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy     uint           `json:"updated_by,omitempty"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy     uint           `json:"-"`
}

func (*Country) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_countries"
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	countries := make([]Country, 0)

	tx := db.Model(Country{}).Order("iso_alpha_code_3 asc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&countries).Error
	if err != nil {
		return nil, err
	}
	p.TotalPage = 1
	p.Items = countries
	return p, nil
}
