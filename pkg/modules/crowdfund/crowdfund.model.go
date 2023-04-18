package crowdfund

import (
	"os"
	"time"

	"crowdfund/pkg/helpers/data"

	"gorm.io/gorm"
)

type Crowdfund struct {
	Id               uint           `json:"id,string" gorm:"primaryKey" validate:"required"`
	Name             string         `json:"name" gorm:"type:varchar(255)"`
	Description      string         `json:"description" gorm:"type:varchar(2000)"`
	ImgBase64        string         `json:"img_base64"`
	CategoryId       uint           `json:"category_id" gorm:"type:int4;"`
	Amount           float32        `json:"amount"`
	ProfitPrecent    float32        `json:"profit_precent"`
	IntroductionText string         `json:"introduction_text"`
	RiskText         string         `json:"risk_text"`
	OrgId            uint           `json:"org_id"`
	OrgName          string         `json:"org_name" gorm:"type:varchar(255)"`
	Status           string         `json:"status" gorm:"type:varchar(255)"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
	CreatedAt        data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy        uint           `json:"-"`
	UpdatedAt        data.LocalTime `json:"-" gorm:"autoUpdateTime"`
	UpdatedBy        uint           `json:"-"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
	Balance          uint           `json:"balance"`
	UsersCnt         uint           `json:"-" gorm:"-"`
}

func (*Crowdfund) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".cwd_crowdfund"
}

type ReqCrowdfundCreate struct {
	Name             string  `json:"name" validate:"required"`
	Description      string  `json:"description"`
	ImgBase64        string  `json:"img_base64" validate:"required"`
	CategoryId       uint    `json:"category_id" validate:"required"`
	Amount           float32 `json:"amount" validate:"required"`
	ProfitPrecent    float32 `json:"profit_precent" validate:"required"`
	IntroductionText string  `json:"introduction_text" validate:"required"`
	RiskText         string  `json:"risk_text" validate:"required"`
	StartDate        string  `json:"start_date" validate:"required"`
	EndDate          string  `json:"end_date" validate:"required"`
}

type ReqCrowdfundId struct {
	Id uint `json:"id,string" validate:"required"`
}
