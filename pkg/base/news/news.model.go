package news

import (
	"os"
	"strings"

	"crowdfund/pkg/base/organization"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type News struct {
	Id        uint           `json:"id,string" gorm:"primaryKey"`
	OrgId     uint           `json:"org_id,string"`
	OrgName   string         `json:"org_name"`
	TypeId    uint           `json:"type_id,string" validate:"required"`
	TypeName  string         `json:"type_name"`
	Title     string         `json:"title" gorm:"type:varchar(2000)" validate:"required"`
	Text      string         `json:"text" validate:"required"`
	Img       string         `json:"img" validate:"required"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*News) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_common_news"
}

func NewsList(c *fiber.Ctx, orgId uint) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	pages := make([]News, 0)
	title := strings.ToLower(c.Query("title"))
	tx := db.Model(News{}).Where("org_id=?", orgId)
	if title != "" {
		tx.Where("lower(title) LIKE ?", "%"+title+"%")
	}

	tx.Count(&totalRow)
	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Order("created_at DESC").Find(&pages).Error
	if err != nil {
		return nil, err
	}
	p.Items = pages
	return p, nil
}

func (p *News) Create() error {
	db := database.DBconn
	p.OrgName = organization.GetName(p.OrgId)
	return db.Create(&p).Error
}

func (p *News) Update() error {
	db := database.DBconn
	return db.Updates(p).Error
}

func (p *News) Delete() error {
	db := database.DBconn
	return db.Delete(p).Error
}
