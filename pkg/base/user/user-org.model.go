package user

import (
	"os"

	"crowdfund/pkg/database"
)

type OrgsResponse struct {
	Id       uint   `json:"id,string"`
	Name     string `json:"name"`
	TypeId   uint   `json:"type_id,string"`
	TypeName string `json:"type_name"`
}

func OrgList(userId int) (*[]OrgsResponse, error) {

	result := make([]OrgsResponse, 0)
	db := database.DBconn
	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_organizations as o").Select("o.id, o.name, o.type_id, o.type_name").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_organization_users as ou ON (o.id = ou.org_id AND o.deleted_at is NULL)").Where("ou.user_id=?", userId)

	err := tx.Find(&result).Error

	return &result, err
}
func OrgListById(userId int, orgId uint) (OrgsResponse, error) {

	result := OrgsResponse{}
	db := database.DBconn
	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_organizations as o").Select("o.id, o.name, o.reg_no, o.type_id, o.type_name").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_organization_users as ou ON (o.id = ou.org_id AND o.deleted_at is NULL)").Where("ou.user_id=? AND ou.org_id=?", userId, orgId)

	err := tx.Find(&result).Error

	return result, err
}
