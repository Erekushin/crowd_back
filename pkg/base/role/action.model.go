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

type Action struct {
	Id          uint           `json:"id,string" gorm:"primaryKey"`
	ModuleId    uint           `json:"module_id,string"`
	ModuleName  string         `json:"module_name" gorm:"type:varchar(120)"`
	PageId      uint           `json:"page_id,string" validate:"required"`
	PageName    string         `json:"page_name" gorm:"type:varchar(120)"`
	Name        string         `json:"name" gorm:"type:varchar(120)" validate:"required"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   data.LocalTime `json:"-" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"created_by,omitempty"`
	UpdatedAt   data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy   uint           `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy   uint           `json:"-"`
}

func (*Action) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbr_actions"
}

func ActionList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	actions := make([]Action, 0)
	name := strings.ToLower(c.Query("name"))
	pageId := convertor.StringToInt(c.Query("page_id"))

	tx := db.Table(os.Getenv("DB_SCHEMA") + ".tbr_actions as a").Select("a.id, a.module_id, a.module_name, a.page_id, a.page_name, a.name, a.description").Joins("LEFT JOIN " + os.Getenv("DB_SCHEMA") + ".tbr_pages as p ON (a.page_id = p.id AND p.deleted_at is NULL)")

	if name != "" {
		tx.Where("lower(a.name) LIKE ?", "%"+name+"%")
	}

	if pageId != 0 {
		tx.Where("a.page_id=?", pageId)
	}

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&actions).Error
	if err != nil {
		return nil, err
	}
	p.Items = actions
	return p, nil
}

func (p *Action) Create() error {
	db := database.DBconn
	page := Page{}
	db.Find(&page, "id=?", p.PageId)

	p.ModuleId = page.ModuleId
	p.ModuleName = page.ModuleName
	p.PageName = page.Name

	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Action) Update() error {
	db := database.DBconn
	page := Page{}
	db.Find(&page, "id=?", p.PageId)

	p.ModuleId = page.ModuleId
	p.ModuleName = page.ModuleName
	p.PageName = page.Name
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Action) Delete() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}
