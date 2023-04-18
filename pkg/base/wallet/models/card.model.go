package account_model

import (
	"os"

	"crowdfund/pkg/helpers/data"

	"gorm.io/gorm"
)

type CardTokens struct {
	Id            uint   `json:"id" gorm:"primaryKey"`
	BankCode      string `json:"-"`
	BankName      string `json:"-"`
	ErrorDesc     string `json:"-"`
	Checksum      string `json:"-"`
	ErrorCode     string `json:"-"`
	Cardholder    string `json:"card_holder"`
	TransactionId string `json:"-"`
	CardNumber    string `json:"card_number"`
	Token         string `json:"-"`
	UserId        uint   `json:"-"`
	Invoice       string `json:"-"`
	Status        uint   `gorm:"type:int4;default:1" json:"-"`
	Bank          Bank   `gorm:"foreignKey:Code;references:BankCode" json:"bank"`

	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy uint           `json:"created_by,omitempty"`
	UpdatedAt data.LocalTime `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedBy uint           `json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy uint           `json:"-"`
}

func (*CardTokens) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_card_tokens"
}

type CardTokenInvoices struct {
	Id            uint `gorm:"primaryKey"`
	TransactionId string
	Checksum      string
	Invoice       string
	UserId        uint
	Amount        float32
	HasToken      uint
	IsConfirm     uint
	RedirectUrl   string
}

func (*CardTokenInvoices) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_card_token_invoices"
}

type CardTokenPayments struct {
	Id            uint
	Amount        float32
	ChargePercent uint
	ChargeAmount  float32
	Hash          string
	DeviceType    string
	ErrorDesc     string
	ErrorCode     string
	Checksum      string
	TransactionId string
	CardNumber    string
	UserId        uint
	CardTokenId   uint
	Type          string
}

func (*CardTokenPayments) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_card_token_payments"
}

type ReqCardAdd struct {
	TransactionId string `json:"invoice" validate:"required"`
	Desc          string `json:"desc"`
	StatusCode    string `json:"status_code"`
}

type ReqCardDelete struct {
	Id uint `json:"id,string" validate:"required"`
}

type ReqCardPay struct {
	Hash          string  `json:"hash" validate:"required"`
	CardTokenId   uint    `json:"card_token_id,string" validate:"required"`
	DeviceType    string  `json:"device_type" validate:"required"`
	Amount        float32 `json:"amount" validate:"required"`
	ChargePercent uint    `json:"charge_percent" validate:"required"`
	ChargeAmount  float32 `json:"charge_amount" validate:"required"`
	CardNo        string  `json:"card_no" validate:"required"`
}

type ResCardDepsitResult struct {
	RefNo string `json:"ref_no"`
}
type ResCardDepsit struct {
	Code    uint                `json:"code"`
	Message string              `json:"message"`
	Result  ResCardDepsitResult `json:"result"`
}
