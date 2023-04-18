package message

import (
	"os"

	"crowdfund/pkg/database"
)

type Message struct {
	Id        uint   `json:"id,string" gorm:"primaryKey"`
	Code      uint   `json:"code" gorm:"unique;type:int4"`
	Name      string `json:"name" gorm:"type:varchar(200)"`
	Path      string `json:"path" gorm:"type:varchar(200)"`
	Method    string `json:"method" gorm:"type:varchar(20);default:POST"`
	IsQuery   uint   `json:"type" gorm:"type:int4;default:0"`
	IsFileLog uint   `json:"is_file_log" gorm:"type:int4;default:0"`
	IsDbLog   uint   `json:"is_db_log" gorm:"type:int4;default:0"`
	IsPublic  uint   `json:"is_public" gorm:"type:int4;default:0"`
	TimeOut   uint   `json:"time_out" gorm:"type:int4;default:10"`
}

func (*Message) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbs_messages"
}

func (m *Message) FindByCode() error {
	db := database.DBconn
	return db.Where("code = ?", m.Code).First(&m).Error
}

func (m *Message) FindByPath() error {
	db := database.DBconn
	return db.Where("path = ? AND method = ?", m.Path, m.Method).First(&m).Error
}
