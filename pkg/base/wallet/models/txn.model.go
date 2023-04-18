package account_model

import (
	"os"

	"crowdfund/pkg/helpers/data"
)

type Transactions struct {
	Id                 int            `json:"-"`
	SrcAccountNo       string         `json:"-"`
	SrcRunningBalance  float32        `json:"-"`
	DestAccountNo      string         `json:"-"`
	DestRunningBalance float32        `json:"-"`
	Description        string         `json:"description"`
	Amount             float32        `json:"amount"`
	TranType           string         `json:"tran_type"`
	PaymentMethod      string         `json:"payment_method"`
	BankCode           string         `json:"bank_code"`
	BankName           string         `json:"bank_name"`
	BankAccount        string         `json:"bank_account"`
	RefNo              string         `json:"ref_no"`
	JournalNo          string         `json:"journal_no"`
	CreatedDate        data.LocalTime `json:"-" gorm:"autoCreateTime"`
}

func (*Transactions) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_transactions"
}

type TxnSignLog struct {
	Id        uint `gorm:"primaryKey"`
	Type      string
	UserId    uint
	Amount    float32
	BankJrNo  string
	AppJrNo   string
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
}

func (*TxnSignLog) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_txn_sign_log"
}

type RespStatement struct {
	Balance float32         `json:"balance"`
	Items   []StatementList `json:"items"`
}

type StatementList struct {
	Id            int            `json:"id"`
	Amount        float32        `json:"amount"`
	TranType      string         `json:"tran_type"`
	JournalNo     string         `json:"journal_no"`
	PaymentMethod string         `json:"payment_method"`
	BankName      string         `json:"bank_name"`
	Description   string         `json:"description"`
	CreatedDate   data.LocalTime `json:"created_date"`
	Date          string         `json:"date"`
	Time          string         `json:"time"`
}
