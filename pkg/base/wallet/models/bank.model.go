package account_model

import "os"

type Bank struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Code string `json:"code" gorm:"varchar(40)"`
	Name string `json:"name" gorm:"varchar(255)"`
	Img  string `json:"img"`
	Logo string `json:"logo"`
}

func (*Bank) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_banks"
}

type UserBankAccounts struct {
	Id            uint   `json:"id" gorm:"primaryKey"`
	UserId        uint   `json:"-"`
	BankId        uint   `json:"-"`
	Bank          Bank   `gorm:"foreignKey:BankId" json:"bank"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

func (*UserBankAccounts) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_user_bank_accounts"
}

type ReqAddBankAccount struct {
	BankId        uint   `json:"bank_id,string" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
}

type ReqDeleteBankAccount struct {
	Id uint `json:"id,string" validate:"required"`
}
