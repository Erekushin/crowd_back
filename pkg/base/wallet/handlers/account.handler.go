package account_handler

import (
	"fmt"
	"os"
	"strings"

	"crowdfund/pkg/base/user"
	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/core"
	"crowdfund/pkg/database"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountHandler struct{}

func (u *AccountHandler) GetDepositBankAccount(c *fiber.Ctx) error {
	return core.Resolve(200, c, core.Response("success", fiber.Map{"account": os.Getenv("WALLET_DEPOSIT_ACCOUNT")}))
}

func (u *AccountHandler) CreateWalletAccount(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	var newAccount account_model.Account

	user := user.ById(userId)
	if user == nil {
		return core.Resolve(400, c, core.Response("User wrong information"))
	}
	newAccount.OwnerId = userId
	newAccount.Name = user.FirstName
	newAccount.Label = user.FirstName
	newAccount.TypeId = 10
	newAccount.TypeName = "Хувь хүн"
	if err := newAccount.Create(); err != nil {
		return core.Resolve(400, c, core.Response("account create err:"+err.Error()))
	}
	return core.Resolve(200, c, core.Response("success"))
}

func (u *AccountHandler) SetDefaultAccount(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	defaultAccountNo := c.Query("account_no")
	if defaultAccountNo == "" {
		return core.Resolve(400, c, core.Response("account_no is required"))
	}

	db := database.DBconn

	if err := db.Model(account_model.Account{}).Where("owner_id=?", userId).Update("is_default", 0).Error; err != nil {
		return core.Resolve(400, c, core.Response("set account default account err:"+err.Error()))
	}

	if err := db.Model(account_model.Account{}).Where("account_no=? AND owner_id=?", defaultAccountNo, userId).Update("is_default", 1).Error; err != nil {
		return core.Resolve(400, c, core.Response("set account default account err:"+err.Error()))
	}

	return core.Resolve(200, c, core.Response("success"))
}

func (u *AccountHandler) GetBalance(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	accounts := make([]account_model.Account, 0)
	db := database.DBconn
	err := db.Where("owner_id = ? AND status = 'A'", userId).Order("is_default DESC").Find(&accounts).Error
	if err != nil {
		return core.Resolve(400, c, core.Response("User wrong information"))
	}

	if len(accounts) == 0 {
		fmt.Println("wallet account create ---1-")
		user := user.ById(userId)
		fmt.Println(user.Id)
		if user.Id != 0 {
			fmt.Println("wallet account create ---2-")
			var account account_model.Account
			account.OwnerId = userId
			account.Name = user.FirstName
			account.Label = user.FirstName
			account.TypeId = 10
			account.TypeName = "Хувь хүн"
			account.IsDefault = 1
			if err := account.Create(); err != nil {
				fmt.Println("wallet account create: ", err.Error())
			} else {
				accounts = append(accounts, account)
			}
		}
	} else {
		fmt.Println("wallet account create hiiihgui")
	}

	return core.Resolve(200, c, core.Response("success", accounts))
}

func (u *AccountHandler) DepositAccount(c *fiber.Ctx) error {

	ip := c.IP()

	if ip == "10.0.0.104" || ip == "127.0.0.1" {
		req := new(account_model.ReqAccountDeposit)
		if err := c.BodyParser(req); err != nil {
			return core.Resolve(400, c, core.Response(err.Error()))
		}

		if errors := core.Validate(*req); errors != nil {
			return core.Resolve(400, c, core.Response("validation error", errors))
		}

		splitDesc := strings.Split(req.Description, "CFU-")

		if len(splitDesc) != 2 {
			return core.Resolve(400, c, core.Response("invalid description"))
		}

		accountNo := splitDesc[1]
		err := AccountDeposit(accountNo, req.RefNo, req.Amount)
		if err != nil {
			return core.Resolve(500, c, core.Response("depo_account error:"+err.Error()))
		}
		return core.Resolve(200, c, core.Response("success"))
	} else {
		return core.Resolve(403, c, core.Response("forbidden"))
	}

}

func (u *AccountHandler) WithDrawAccount(c *fiber.Ctx) error {
	req := new(account_model.ReqAccountWithDraw)
	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if req.Amount < 500 {
		return core.Resolve(400, c, core.Response("invalid amount"))
	}

	var fee float32 = 200
	var transfer_money float32 = 0
	var bankAccount account_model.UserBankAccounts

	userId := oauth.GetSessionUserId(c)
	db := database.DBconn

	err := db.Where("id = ? AND user_id = ?", req.UserBankAccountId, userId).Preload("Bank").First(&bankAccount).Error
	if err != nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	if bankAccount.Bank.Code == "040000" {
		fee = 100
	}

	transfer_money = req.Amount - fee

	user := user.ById(userId)
	if user == nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	account, _ := account_model.GetAccount(req.AccountNumber)
	if userId != account.OwnerId {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	if account.Balance < float32(req.Amount) {
		return core.Resolve(400, c, core.Response("blance not enough"))
	}

	journal_no := uuid.New().String()
	tdbWithdrawAccount := "2000000005"
	reqWithdraw := new(account_model.ReqCgWithdraw)
	reqWithdraw.Amount = transfer_money
	reqWithdraw.Description = "Шилжүүлэг"
	reqWithdraw.SrcAcc = os.Getenv("CG_TDB_WITHDRAW_SRC_ACCOUNT")
	reqWithdraw.DestAcc = bankAccount.AccountNumber
	reqWithdraw.DestName = bankAccount.AccountName
	reqWithdraw.DestBank = bankAccount.Bank.Code

	if err := MakeTransaction(req.Amount, journal_no, account.AccountNo, tdbWithdrawAccount, "WITHDRAW", "bank", bankAccount.Bank.Code, bankAccount.Bank.Name, bankAccount.AccountNumber, ""); err != nil {
		return core.Resolve(500, c, core.Response("withdraw transaction error: "+err.Error()))
	}

	err, ref_no := bankTransfer(*reqWithdraw)
	if err != nil {
		fmt.Println("\n\n\nbankTransfer error: " + err.Error())
	}

	txnSignLog := new(account_model.TxnSignLog)
	txnSignLog.Type = "-"
	txnSignLog.Amount = float32(req.Amount)
	txnSignLog.BankJrNo = ref_no
	txnSignLog.AppJrNo = journal_no
	txnSignLog.UserId = userId
	db.Save(&txnSignLog)

	return core.Resolve(200, c, core.Response("success"))
}

func (u *AccountHandler) Send(c *fiber.Ctx) error {
	req := new(account_model.ReqSendAmount)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if req.Amount < 1 {
		return core.Resolve(400, c, core.Response("invalid amount"))
	}

	var err error
	srcUserId := oauth.GetSessionUserId(c)
	srcAccount, err := account_model.CheckAccountOwner(srcUserId, req.SrcAccountNo)
	if err != nil {
		return core.Resolve(400, c, core.Response("invalid src account"))
	}

	if srcAccount.Balance < float32(req.Amount) {
		return core.Resolve(400, c, core.Response("balance not enough"))
	}

	destAccount, err := account_model.GetDefaultAccount(req.DestUserId)
	if err != nil {
		return core.Resolve(400, c, core.Response("invalid dest account"))
	}

	tdbMirrorAccount := "2000000003"
	journal_no := uuid.New().String()

	if err := MakeTransaction(req.Amount, journal_no, srcAccount.AccountNo, tdbMirrorAccount, req.Description, "", "", "", "", ""); err != nil {
		return core.Resolve(400, c, core.Response("transaction error1:"+err.Error()))
	}

	if err := MakeTransaction(req.Amount, journal_no, tdbMirrorAccount, destAccount.AccountNo, req.Description, "", "", "", "", ""); err != nil {
		return core.Resolve(400, c, core.Response("transaction error2:"+err.Error()))
	}

	// notif_text := currentUser.FirstName + "-с мөнгө ирлээ"
	// go client.SendNotif(req.DestUserId, notif_text, notif_text, "txn", strconv.Itoa(transaction_id))

	return core.Resolve(200, c, core.Response("success"))
}

func (u *AccountHandler) StatementList(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	accountNo := c.Query("account_no")

	var (
		account account_model.Account
		err     error
	)

	if accountNo == "" {
		account, _ = account_model.GetDefaultAccount(userId)
		if err != nil {
			return core.Resolve(400, c, core.Response("invalid account"))
		}
	} else {
		account, err = account_model.CheckAccountOwner(userId, accountNo)
		if err != nil {
			return core.Resolve(400, c, core.Response("invalid account"))
		}
	}

	db := database.DBconn

	var statement account_model.RespStatement
	transactions := make([]account_model.StatementList, 0)

	currentUser := user.ById(userId)
	if currentUser == nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	db.Model(&account_model.Transactions{}).Select("*, to_char(created_date, 'YYYY-MM-DD') as date, to_char(created_date, 'HH24:MI') as time").Order("created_date desc").Where("src_account_no = ? ", account.AccountNo).Limit(50).Scan(&transactions)
	db.Where("owner_id = ?", userId).First(&account)

	statement.Balance = account.Balance
	statement.Items = transactions
	return core.Resolve(200, c, core.Response("success", statement))
}
