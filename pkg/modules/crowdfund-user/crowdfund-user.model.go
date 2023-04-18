package crowdfunduser

import (
	"crowdfund/pkg/helpers/data"
	"os"

	"gorm.io/gorm"
)

type CrowdfundUser struct {
	Id          uint           `json:"id,string" gorm:"primaryKey" validate:"required"`
	CrowdfundId uint           `json:"fund_id" validate:"required"`
	UserId      uint           `json:"user_id" validate:"required"`
	Amount      uint           `json:"amount" validate:"required"`
	CreatedAt   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"-"`
	UpdatedAt   data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy   uint           `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy   uint           `json:"-"`
}

func (*CrowdfundUser) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".cwd_crowdfund_user"
}

type ReqCrowdfundUserCreate struct {
	CrowdfundId uint `json:"fund_id" validate:"required"`
	UserId      uint `json:"user_id" validate:"required"`
	Amount      uint `json:"amount" validate:"required"`
}
type ReqCrowdfundUserList struct {
	CrowdfundId uint `json:"fund_id" validate:"required"`
}

type CrowdfundUsers struct {
	UserId      uint   `json:"user_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	RegNo       string `json:"reg_no" validate:"required"`
	CrowdfundId uint   `json:"crowdfund_id" validate:"required"`
	Amount      uint   `json:"amount" validate:"required"`
}

type UserCrowdfund struct {
	Id             uint           `json:"id,string" validate:"required"`
	Name           string         `json:"name" gorm:"type:varchar(255)"`
	Description    string         `json:"description" gorm:"type:varchar(2000)"`
	ImgBase64      string         `json:"img_base64"`
	CategoryId     uint           `json:"category_id" gorm:"type:int4;"`
	Amount         float32        `json:"amount"`
	ProfitPrecent  float32        `json:"profit_precent"`
	OrgId          uint           `json:"org_id"`
	Status         string         `json:"status" gorm:"type:varchar(255)"`
	StartDate      data.LocalDate `json:"start_date"`
	EndDate        data.LocalDate `json:"end_date"`
	CreatedAt      data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy      uint           `json:"-"`
	UpdatedAt      data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy      uint           `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy      uint           `json:"-"`
	InvestedAmount uint           `json:"invested_amount"`
}
