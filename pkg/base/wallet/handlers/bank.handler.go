package account_handler

import (
	"crowdfund/pkg/base/user"
	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/core"
	"crowdfund/pkg/database"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type BankHandler struct{}

func (u *BankHandler) BankList(c *fiber.Ctx) error {
	var banks []account_model.Bank
	db := database.DBconn
	db.Order("id asc").Find(&banks)
	return core.Resolve(200, c, core.Response("success", banks))
}

func (u *BankHandler) AddUserBankAccount(c *fiber.Ctx) error {
	req := new(account_model.ReqAddBankAccount)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	var bank_account account_model.UserBankAccounts
	userId := oauth.GetSessionUserId(c)
	currentUser := user.ById(userId)
	if currentUser == nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	bank_account.UserId = oauth.GetSessionUserId(c)
	bank_account.BankId = req.BankId
	bank_account.AccountName = currentUser.FirstName
	bank_account.AccountNumber = req.AccountNumber

	db := database.DBconn
	if err := db.Create(&bank_account).Error; err != nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}
	return core.Resolve(200, c, core.Response("success"))
}

func (u *BankHandler) GetUserBankAccounts(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	db := database.DBconn
	var banks_accounts []account_model.UserBankAccounts
	db.Preload("Bank").Find(&banks_accounts, "user_id = ?", userId)

	return core.Resolve(200, c, core.Response("success", banks_accounts))
}

func (u *BankHandler) DeleteUserBankAccount(c *fiber.Ctx) error {
	req := new(account_model.ReqDeleteBankAccount)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	bank_account := account_model.UserBankAccounts{Id: req.Id, UserId: oauth.GetSessionUserId(c)}

	db := database.DBconn
	if err := db.Delete(bank_account).Error; err != nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}
	return core.Resolve(200, c, core.Response("success"))
}
