package orgtypeaction

import (
	"os"

	"crowdfund/pkg/base/role"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"

	"github.com/gofiber/fiber/v2"
)

type OrgTypeAction struct {
	TypeId    uint           `json:"type_id,string" validate:"required"`
	ActionId  uint           `json:"action_id,string" validate:"required"`
	CreatedAt data.LocalTime `json:"-" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"-"`
}

type OrgTypeActionResult struct {
	TypeId   uint `json:"type_id,string"`
	ActionId uint `json:"action_id,string"`
}

func (*OrgTypeAction) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbr_org_type_actions"
}

func OrgTypeActionList(c *fiber.Ctx, typeId int) []role.Action {
	result := make([]role.Action, 0)
	db := database.DBconn

	db.Table(os.Getenv("DB_SCHEMA")+".tbr_actions as a").Select("a.id, a.module_id, a.module_name, a.page_id, a.page_name, a.name, a.description").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbr_pages as p ON (a.page_id = p.id AND p.deleted_at is NULL)").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbr_org_type_actions as ota ON (ota.action_id = a.id AND a.deleted_at is NULL)").Where("ota.type_id=?", typeId).Find(&result)

	return result
}

func (r *OrgTypeAction) Add() error {
	db := database.DBconn
	return db.Create(&r).Error
}

func (r *OrgTypeAction) Remove() error {
	db := database.DBconn
	return db.Where("type_id=? AND action_id=?", r.TypeId, r.ActionId).Delete(&r).Error
}
