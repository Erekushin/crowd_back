package terminal

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Terminal struct {
	Id          uint           `json:"id,string"`
	Name        string         `json:"name"`
	SerialNo    string         `json:"serial_no" gorm:"unique;not null"`
	GroupId     uint           `json:"group_id,string"`
	ConfigJson  datatypes.JSON `json:"config_json"`
	Description string         `json:"description"`
	TellerId    uint           `json:"teller_id"`
	LocationId  uint           `json:"location_id"`
	CreatedAt   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"created_by,omitempty"`
	UpdatedAt   data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy   uint           `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy   uint           `json:"-"`
}

func GetTerminalTableName(alias string) string {
	return os.Getenv("DB_SCHEMA") + ".tbd_terminals as " + alias
}

func (*Terminal) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_terminals"
}

func (p *Terminal) Create() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	terminals := make([]Terminal, 0)
	serial_no := strings.ToLower(c.Query("serial_no"))

	tx := db.Model(Terminal{})
	if serial_no != "" {
		tx.Where("lower(serial_no) LIKE ? OR lower(name) LIKE ?", "%"+serial_no+"%", "%"+serial_no+"%")
	}
	tx.Order("created_at desc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&terminals).Error
	if err != nil {
		return nil, err
	}
	p.Items = terminals
	return p, nil
}

func (p *Terminal) Update() error {
	db := database.DBconn
	if err := db.Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Terminal) Delete() error {
	db := database.DBconn
	if err := db.Delete(p).Error; err != nil {
		return err
	}
	return nil
}
