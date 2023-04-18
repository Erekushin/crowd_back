package role

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Page struct {
	Id         uint           `json:"id,string" gorm:"primaryKey"`
	ModuleId   uint           `json:"module_id,string" validate:"required"`
	ModuleName string         `json:"module_name" gorm:"type:varchar(120)"`
	Code       string         `json:"code" gorm:"type:varchar(30)" validate:"required"`
	Name       string         `json:"name" gorm:"type:varchar(120)" validate:"required"`
	Path       string         `json:"path" gorm:"type:varchar(30)" validate:"required"`
	Icon       string         `json:"icon" gorm:"type:varchar(30)"  validate:"required"`
	Sequence   uint           `json:"sequence" gorm:"type:int4;default:999"`
	GroupCode  string         `json:"group_code" gorm:"type:varchar(30)" validate:"required"`
	GroupName  string         `json:"group_name" gorm:"type:varchar(120)" validate:"required"`
	CreatedAt  data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy  uint           `json:"created_by,omitempty"`
	UpdatedAt  data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy  uint           `json:"updated_by,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy  uint           `json:"-"`
}

func (*Page) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbr_pages"
}

func PageList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	pages := make([]Page, 0)
	name := strings.ToLower(c.Query("name"))
	moduleId := convertor.StringToInt(c.Query("module_id"))
	tx := db.Model(Page{})
	if name != "" {
		tx.Where("lower(name) LIKE ?", "%"+name+"%")
	}

	if moduleId != 0 {
		tx.Where("module_id=?", moduleId)
	}
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Order("sequence").Find(&pages).Error
	if err != nil {
		return nil, err
	}
	p.Items = pages
	return p, nil
}

func (p *Page) Create() error {
	db := database.DBconn

	module := Module{}
	db.Find(&module, "id=?", p.ModuleId)
	p.ModuleName = module.Name

	if err := db.Create(&p).Error; err != nil {
		return err
	}

	action := Action{
		ModuleId:   p.ModuleId,
		ModuleName: p.ModuleName,
		PageId:     p.Id,
		PageName:   p.Name,
		Name:       p.Name + " харах",
	}

	if err := db.Create(&action).Error; err != nil {
		return nil
	}

	return nil
}

func (p *Page) Update() error {
	db := database.DBconn

	module := Module{}
	db.Find(&module, "id=?", p.ModuleId)
	p.ModuleName = module.Name

	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Page) Delete() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}

	if err := db.Where("page_id=?", p.Id).Delete(&Action{}).Error; err != nil {
		return err
	}

	return nil
}
