package organization

import (
	"os"
	"strings"

	"crowdfund/pkg/base/country"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Organization struct {
	Id            uint           `json:"id,string" gorm:"primaryKey"`
	RegNo         string         `json:"reg_no" gorm:"type:varchar(7)"`
	Name          string         `json:"name" gorm:"type:varchar(255)" validate:"required"`
	ShortName     string         `json:"short_name" gorm:"type:varchar(255)"`
	RootAccount   uint           `json:"root_account,string"`
	PhoneNo       string         `json:"phone_no" gorm:"type:varchar(20)"`
	Email         string         `json:"email" gorm:"type:varchar(50)"`
	LogoImage     string         `json:"logo_image"`
	AimagCode     string         `json:"aimag_code" gorm:"type:varchar(5)"`
	AimagName     string         `json:"aimag_name" gorm:"type:varchar(255)"`
	SumCode       string         `json:"sum_code" gorm:"type:varchar(5)"`
	SumName       string         `json:"sum_name" gorm:"type:varchar(255)"`
	BagCode       string         `json:"bag_code" gorm:"type:varchar(5)"`
	BagName       string         `json:"bag_name" gorm:"type:varchar(255)"`
	Address       string         `json:"address" gorm:"type:varchar(600)"`
	TypeId        uint           `json:"type_id,string" validate:"required"`
	TypeName      string         `json:"type_name" gorm:"type:varchar(200)"`
	CountryCode   string         `json:"country_code" gorm:"type:varchar(3)"`
	CountryName   string         `json:"country_name"`
	CountryNameEn string         `json:"country_name_en"`
	CreatedAt     data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy     uint           `json:"created_by,omitempty"`
	UpdatedAt     data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy     uint           `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (*Organization) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organizations"
}

func GetName(id uint) string {
	db := database.DBconn
	var ct Organization
	if err := db.First(&ct, id).Error; err != nil {
		return "Тодорхойгүй"
	}
	return ct.Name
}

func (p *Organization) Create() error {
	db := database.DBconn
	if p.ShortName == "" {
		p.ShortName = p.Name
	}

	if p.CountryCode != "" {
		c := new(country.Country)

		if err := db.First(&c, "iso_alpha_code_3=?", strings.ToUpper(p.CountryCode)); err != nil {
			p.CountryName = c.CommonName
			p.CountryNameEn = c.EnName
		}
	} else {
		p.CountryCode = "MNG"
		p.CountryName = "Монгол"
		p.CountryNameEn = "Mongolia"
	}

	p.TypeName = GetTypeName(p.TypeId)
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	organizations := make([]Organization, 0)
	searchText := strings.ToLower(c.Query("search_text"))
	typeId := convertor.StringToInt(c.Query("type_id"))

	tx := db.Model(Organization{})
	if searchText != "" {
		tx.Where("reg_no LIKE ? OR lower(name) LIKE ?", "%"+searchText+"%", "%"+searchText+"%")
	}

	if typeId != 0 {
		tx.Where("type_id = ?", typeId)
	}

	tx.Order("created_at desc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	p.Items = organizations
	return p, nil
}

func (p *Organization) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Organization) Remove() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}

func FindById(orgId uint) *Organization {
	db := database.DBconn
	var org Organization
	if err := db.First(&org, "id =?", orgId).Error; err != nil {
		return nil
	}

	return &org
}
