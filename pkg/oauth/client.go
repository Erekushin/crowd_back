package oauth

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type ClientInfo struct {
	RequestUserId     uint
	RequestOrgId      uint
	RequestTerminalId uint
}

type Client struct {
	Id           uint   `gorm:"primaryKey"`
	ClientId     string `gorm:"type:varchar(50)"`
	ClientSecret string `gorm:"type:varchar(50)"`
	Name         string `gorm:"type:varchar(80)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (*Client) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".oauth_clients2"
}
