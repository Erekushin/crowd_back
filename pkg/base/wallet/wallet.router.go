package wallet

import (
	account_handler "crowdfund/pkg/base/wallet/handlers"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var accountHandler account_handler.AccountHandler
	var bankHandler account_handler.BankHandler
	var cardHandler account_handler.CardService
	var invoiceHandler account_handler.InvoiceHandler

	walletApi := app.Group("wallet")
	walletApi.Post("deposit", accountHandler.DepositAccount)
	walletApi.Get("deposit/account", accountHandler.GetDepositBankAccount)
	walletApi.Post("withdraw", oauth.TokenMiddleware, accountHandler.WithDrawAccount)
	walletApi.Post("send", oauth.TokenMiddleware, accountHandler.Send)

	account := walletApi.Group("account", oauth.TokenMiddleware)
	account.Get("balance", accountHandler.GetBalance)
	account.Post("", accountHandler.CreateWalletAccount)
	account.Post("default", accountHandler.SetDefaultAccount)
	account.Get("statement", accountHandler.StatementList)

	bankApi := walletApi.Group("bank", oauth.TokenMiddleware)
	bankApi.Get("", bankHandler.BankList)
	bankApi.Post("account", bankHandler.AddUserBankAccount)
	bankApi.Get("account", bankHandler.GetUserBankAccounts)
	bankApi.Post("account/delete", bankHandler.DeleteUserBankAccount)

	cardApi := walletApi.Group("card")
	cardApi.Get("", oauth.TokenMiddleware, cardHandler.GetCardList)
	cardApi.Post("delete", oauth.TokenMiddleware, cardHandler.CardDelete)
	cardApi.Post("invoice", oauth.TokenMiddleware, cardHandler.CreateCardInvoice)
	cardApi.Post("deposit", oauth.TokenMiddleware, cardHandler.CardDeposit)

	invoiceApi := walletApi.Group("invoice", oauth.TokenMiddleware)
	invoiceApi.Get("", invoiceHandler.GetList)
	invoiceApi.Post("", invoiceHandler.CreateInvoice)
	invoiceApi.Post("pay", invoiceHandler.PayInvoice)
	invoiceApi.Post("cancel", invoiceHandler.CancelInvoice)
}
