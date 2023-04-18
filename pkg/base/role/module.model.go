package role

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Module struct {
	Id          uint           `json:"id,string" gorm:"primaryKey"`
	Code        string         `json:"code" gorm:"type:varchar(30)" validate:"required"`
	Name        string         `json:"name" gorm:"type:varchar(120)" validate:"required"`
	Icon        string         `json:"icon" gorm:"type:varchar(120)" validate:"required"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	Sequence    int            `json:"sequence" gorm:"default:999"`
	CreatedAt   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"created_by,omitempty"`
	UpdatedAt   data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy   uint           `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy   uint           `json:"-"`
}

func (*Module) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbr_modules"
}

func ModuleList(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	modules := make([]Module, 0)
	name := strings.ToLower(c.Query("name"))

	tx := db.Model(Module{}).Where("lower(name) LIKE ?", "%"+name+"%")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&modules).Error
	if err != nil {
		return nil, err
	}
	p.Items = modules
	return p, nil
}

func (p *Module) Create() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Module) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Module) Delete() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}

	if err := db.Model(&Page{}).Where("module_id=?", p.Id).Updates(Page{ModuleId: 0, ModuleName: ""}).Error; err != nil {
		return err
	}

	return nil
}
