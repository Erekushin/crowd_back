package otp

import (
	"os"
	"time"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"
)

type OtpCode struct {
	Id        uint
	Identity  string
	Code      uint
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
}

func (*OtpCode) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbs_otp_codes"
}

func (p *OtpCode) CheckOtp(identity string, maxAllowedTime time.Time) int64 {
	db := database.DBconn
	return db.Where("identity = ? AND created_at > ?", identity, maxAllowedTime).Find(&p).RowsAffected
}

func (p *OtpCode) GetLastOtp(identity string) error {
	db := database.DBconn
	if err := db.Where("identity = ?", identity).Last(&p).Error; err != nil {
		return err
	}
	return nil
}

func (p *OtpCode) Save() error {
	db := database.DBconn
	if err := db.Create(p).Error; err != nil {
		return err
	}
	return nil
}
