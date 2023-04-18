package migration

import (
	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/database"
	"crowdfund/pkg/modules/crowdfund"
	crowdfunduser "crowdfund/pkg/modules/crowdfund-user"
)

func Run() {
	db := database.DBconn
	// db.AutoMigrate(message.Message{}, core.Log{})
	// db.AutoMigrate(oauth.Client{}, oauth.Session{})
	db.AutoMigrate(account_model.Account{}, account_model.AccountType{}, account_model.Bank{}, account_model.UserBankAccounts{})
	// db.AutoMigrate(account_model.TxnSignLog{}, account_model.Transactions{})
	// db.AutoMigrate(account_model.CardTokens{}, account_model.CardTokenInvoices{}, account_model.CardTokenPayments{})
	// db.AutoMigrate(account_model.Invoice{})
	// db.AutoMigrate(user.User{})
	// db.AutoMigrate(vehicle.Vehicle{})
	// db.AutoMigrate(terminal.Terminal{})
	// db.AutoMigrate(location.Location{})
	// db.AutoMigrate(organization.Organization{}, organization.OrgType{})
	// db.AutoMigrate(lang.Lang{}, lang.Key{}, lang.Translation{})
	// db.AutoMigrate(role.Role{}, role.Module{}, role.Page{}, role.Action{})
	// db.AutoMigrate(document.Document{}, document.Type{}, document.Category{})
	// db.AutoMigrate(organization.OrgUser{})
	// db.AutoMigrate(otp.OtpCode{})
	// db.AutoMigrate(orgtypeaction.OrgTypeAction{})
	// db.AutoMigrate(country.Country{})
	// db.AutoMigrate(news.News{}, news.NewsType{})
	db.AutoMigrate(crowdfund.Crowdfund{})
	// db.AutoMigrate(category.Category{})
	// db.AutoMigrate(crowdfund.Crowdfund{})
	// db.AutoMigrate(category.Category{})
	db.AutoMigrate(crowdfunduser.CrowdfundUser{})
}
