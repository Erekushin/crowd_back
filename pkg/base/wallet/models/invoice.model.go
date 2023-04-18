package account_model

import (
	"os"

	"crowdfund/pkg/helpers/data"

	"gorm.io/datatypes"
)

type Invoice struct {
	Id            uint           `json:"id" gorm:"primaryKey"`
	SrcUserId     uint           `json:"-"`
	SrcAccountNo  string         `json:"-"`
	DestUserId    uint           `json:"-"`
	DestAccountNo string         `json:"-"`
	Amount        float32        `json:"total_amount"`
	Status        string         `json:"status"` //NEW,PAID,DECLINED
	Type          string         `json:"-"`      //USER_TO_USER, ORG_TO_USER
	Description   string         `json:"description"`
	BpInvoiceBody datatypes.JSON `json:"pay_body"`
	DueDate       data.LocalTime `json:"due_date"`
	CreatedDate   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	ListType      string         `gorm:"-" json:"list_type"`
}

func (*Invoice) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_invoices"
}

type ReqInvoiceCreate struct {
	DestUserId    uint           `json:"dest_user_id,string" validate:"required"`
	Amount        float32        `json:"amount" validate:"required"`
	Description   string         `json:"description"`
	Type          string         `json:"type"`
	BpInvoiceBody datatypes.JSON `json:"body"`
}

type ReqInvoicePay struct {
	Id      int    `json:"id,string" validate:"required"`
	PinCode string `json:"pin_code"`
}

type ReqInvoiceCancel struct {
	Id int `json:"id,string" validate:"required"`
}

type RespInvoiceList struct {
	Id               uint           `json:"id"`
	DestUserId       uint           `json:"-"`
	SrcName          string         `json:"src_name"`
	SrcProfileImage  string         `json:"src_profile_image"`
	DestName         string         `json:"dest_name"`
	DestProfileImage string         `json:"dest_profile_image"`
	Amount           float32        `json:"total_amount"`
	Status           string         `json:"status"`
	Description      string         `json:"description"`
	DueDate          data.LocalTime `json:"due_date"`
	CreatedDate      data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	ListType         string         `gorm:"-" json:"list_type"`
}
