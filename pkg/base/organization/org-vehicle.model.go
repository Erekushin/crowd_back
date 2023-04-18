package organization

import (
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
)

type OrgVehicle struct {
	OrgId     uint           `json:"org_id,string" validate:"required"`
	VehicleId uint           `json:"vehicle_id,string" validate:"required"`
	CreatedAt data.LocalTime `gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
}

func (*OrgVehicle) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbd_organization_vehicles"
}

type VehiclesResponse struct {
	Id            uint           `json:"id,string"`
	PlateNo       string         `json:"plate_no"`
	VehicleTypeId uint           `json:"vehicle_type_id"`
	ColorName     string         `json:"color_name"`
	MaxLoad       uint           `json:"max_load"`
	PurposeId     uint           `json:"purpose_id"`
	PurposeName   string         `json:"purpose_name"`
	SeatsCount    uint           `json:"seats_count"`
	CabinNo       string         `json:"cabin_no"`
	Capacity      uint           `json:"capacity"`
	ModelName     string         `json:"model_name"`
	MarkName      string         `json:"mark_name"`
	OwnerName     string         `json:"owner_name"`
	IsHybrid      uint           `json:"is_hybrid"`
	AtutId        uint           `json:"atut_id"`
	ClassName     string         `json:"class_name"`
	CreatedAt     data.LocalTime `json:"created_date"`
	CreatedBy     uint           `json:"created_by"`
}

func VehicleList(c *fiber.Ctx, orgId int) (*data.Pagination, error) {
	var (
		err      error
		totalRow int64
	)
	result := make([]VehiclesResponse, 0)
	db := database.DBconn
	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_vehicles as v").Select("v.*, ov.created_at, ov.created_by").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_organization_vehicles as ov ON (v.id = ov.vehicle_id AND v.deleted_at is NULL)").Where("ov.org_id=?", orgId)
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err = tx.Offset(p.Offset).Limit(p.PageSize).Find(&result).Error
	if err != nil {
		return nil, err
	}
	p.Items = result
	return p, nil
}

func (r *OrgVehicle) Add() error {
	db := database.DBconn
	if err := db.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrgVehicle) Remove() error {
	db := database.DBconn
	if err := db.Where("org_id=?", r.OrgId).Where("vehicle_id=?", r.VehicleId).Delete(&r).Error; err != nil {
		return err
	}

	return nil
}
