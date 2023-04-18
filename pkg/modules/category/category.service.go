package category

import (
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
)

func (*Category) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".cwd_category"
}

func (c *Category) Create() (*Category, error) {
	db := database.DBconn
	return c, db.Create(&c).Error
}

func (c *Category) Update() error {
	db := database.DBconn
	return db.Where("id=?", c.Id).Updates(c).Error
}

func (c *Category) Delete() error {
	db := database.DBconn
	return db.Delete(c).Error
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	crfund := make([]Category, 0)

	tx := db.Model(Category{})

	tx.Order("created_at asc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&crfund).Error
	if err != nil {
		return nil, err
	}
	p.Items = crfund
	return p, nil
}
