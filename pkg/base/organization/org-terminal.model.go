package organization

import (
	"os"

	"crowdfund/pkg/base/terminal"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type OrgTerminal struct {
	OrgId      uint           `json:"org_id,string" validate:"required"`
	TerminalId uint           `json:"terminal_id,string" validate:"required"`
	CreatedAt  data.LocalTime `gorm:"autoCreateTime"`
	CreatedBy  uint           `json:"created_by,omitempty"`
}

func GetOrgTerminalTableName(alias string) string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organization_terminals as " + alias
}

func (*OrgTerminal) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organization_terminals"
}

type TerminalsResponse struct {
	Id          uint           `json:"id,string"`
	Name        string         `json:"name"`
	SerialNo    string         `json:"serial_no"`
	GroupId     uint           `json:"group_id,string"`
	ConfigJson  datatypes.JSON `json:"config_json"`
	Description string         `json:"description"`
	TellerId    uint           `json:"teller_id"`
	LocationId  uint           `json:"location_id"`
	CreatedAt   data.LocalTime `json:"created_at"`
	CreatedBy   uint           `json:"created_by"`
}

func TerminalList(c *fiber.Ctx, orgId int) (*data.Pagination, error) {
	var (
		err      error
		totalRow int64
	)
	result := make([]TerminalsResponse, 0)
	db := database.DBconn
	tx := db.Table(terminal.GetTerminalTableName("t")).Select("t.*, ot.created_at, ot.created_by").Joins("LEFT JOIN "+GetOrgTerminalTableName("ot")+" ON (t.id = ot.terminal_id AND t.deleted_at is NULL)").Where("ot.org_id=?", orgId)
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err = tx.Offset(p.Offset).Limit(p.PageSize).Find(&result).Error
	if err != nil {
		return nil, err
	}
	p.Items = result
	return p, nil
}

func (r *OrgTerminal) Create() error {
	db := database.DBconn
	if err := db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrgTerminal) Remove() error {
	db := database.DBconn
	if err := db.Where("org_id=?", r.OrgId).Where("terminal_id=?", r.TerminalId).Delete(&r).Error; err != nil {
		return err
	}

	return nil
}
