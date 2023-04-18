package account_model

import (
	"fmt"
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"
)

type Account struct {
	Id        int            `json:"-" gorm:"primaryKey"`
	AccountNo string         `json:"account_no" gorm:"uniqueIndex;type:varchar(20)"`
	OwnerId   uint           `json:"-"`
	Balance   float32        `json:"balance" gorm:"default:0"`
	Status    string         `json:"-" gorm:"type:varchar(20);default:A"`
	Name      string         `json:"name" gorm:"type:varchar(255)"`
	Label     string         `json:"label" gorm:"type:varchar(255)"`
	TypeId    int            `json:"type_id" gorm:"type:int4;default:10"`
	TypeName  string         `json:"type_name" gorm:"type:varchar(255);default:'Хувь хүн'"`
	IsDefault uint           `json:"is_default" gorm:"type:int4;default:0"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
}

func (*Account) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_accounts"
}

type AccountType struct {
	Id        int            `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255)"`
	CreatedAt data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
}

func (*AccountType) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbw_accounts_types"
}

func GetAccount(accountNo string) (c Account, e error) {
	db := database.DBconn
	e = db.First(&c, "account_no=?", accountNo).Error
	return c, e
}

func GetDefaultAccount(userId uint) (c Account, e error) {
	db := database.DBconn
	e = db.First(&c, "owner_id=? AND is_default=1", userId).Error
	return c, e
}

func CheckAccountOwner(userId uint, accountNo string) (c Account, e error) {
	db := database.DBconn
	e = db.First(&c, "owner_id=? AND account_no=?", userId, accountNo).Error
	return c, e
}

func (p *Account) Create() error {
	db := database.DBconn
	var id uint
	db.Raw("SELECT nextval('" + os.Getenv("DB_SCHEMA") + ".tbw_account_no_seq')").Scan(&id)
	p.AccountNo = fmt.Sprintf("%d%08d", p.TypeId, id)
	return db.Create(&p).Error
}

type ReqAccountDeposit struct {
	Description   string  `json:"description" validate:"required"`
	Amount        float32 `json:"amount" validate:"required"`
	RefNo         string  `json:"ref_no" validate:"required"`
	PaymentMethod string  `json:"payment_method"` // card tokeninaze or tdb shiljiileg
	BankCode      string  `json:"bank_code"`      // odoo ashiglahgui bga
	BankName      string  `json:"bank_name"`      // odoo ashiglahgui bga
}

type ReqAccountWithDraw struct {
	UserBankAccountId int     `json:"bank_account_id,string" validate:"required"`
	AccountNumber     string  `json:"account_number" validate:"required"`
	Amount            float32 `json:"amount" validate:"required"`
}

type ReqSendAmount struct {
	DestUserId   uint    `json:"dest_user_id,string" validate:"required"`
	SrcAccountNo string  `json:"src_account_no" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
	Description  string  `json:"description"`
}
