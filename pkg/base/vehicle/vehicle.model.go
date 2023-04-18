package vehicle

import (
	"os"
	"strings"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Vehicle struct {
	Id            uint           `json:"id,string"`
	PlateNo       string         `json:"plate_no" gorm:"unique;not null"`
	VehicleTypeId uint           `json:"vehicle_type_id"`
	ColorName     string         `json:"color_name"`
	MaxLoad       uint           `json:"max_load"`
	PurposeId     uint           `json:"purpose_id" gorm:"type:smallint"`
	PurposeName   string         `json:"purpose_name"`
	SeatsCount    uint           `json:"seats_count"`
	CabinNo       string         `json:"cabin_no"`
	Capacity      uint           `json:"capacity"`
	ModelName     string         `json:"model_name"`
	MarkName      string         `json:"mark_name"`
	OwnerName     string         `json:"owner_name"`
	IsHybrid      uint           `json:"is_hybrid" gorm:"type:smallint"`
	AtutId        uint           `json:"atut_id"`
	ClassName     string         `json:"class_name"`
	CreatedAt     data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy     uint           `json:"created_by,omitempty"`
	UpdatedAt     data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy     uint           `json:"updated_by,omitempty"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy     uint           `json:"-"`
}

func (*Vehicle) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_vehicles"
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	vehicles := make([]Vehicle, 0)
	plate_no := strings.ToLower(c.Query("plate_no"))

	tx := db.Model(Vehicle{})
	if plate_no != "" {
		tx.Where("plate_no = ?", plate_no)
	}
	tx.Order("created_at desc")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&vehicles).Error
	if err != nil {
		return nil, err
	}
	p.Items = vehicles
	return p, nil
}

func (v *Vehicle) Find() error {
	db := database.DBconn
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plate_no"}},
		UpdateAll: true,
	}).Create(&v).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *Vehicle) BeforeSave(tx *gorm.DB) (err error) {
	v.PlateNo = strings.ToLower(v.PlateNo)
	v.ColorName = strings.ToLower(v.ColorName)
	v.PurposeName = strings.ToLower(v.PurposeName)
	v.CabinNo = strings.ToLower(v.CabinNo)
	v.ModelName = strings.ToLower(v.ModelName)
	v.MarkName = strings.ToLower(v.MarkName)
	v.OwnerName = strings.ToLower(v.OwnerName)
	return
}
