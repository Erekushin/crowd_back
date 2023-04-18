package account_handler

import (
	"os"
	"time"

	"crowdfund/pkg/base/user"
	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/core"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvoiceHandler struct{}

func (u *InvoiceHandler) GetList(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)

	db := database.DBconn

	invoices := make([]account_model.RespInvoiceList, 0)

	db.Table(os.Getenv("DB_SCHEMA")+".tbw_invoices i").Select("i.*, su.last_name || ' ' || su.first_name as src_name, su.profile_image src_profile_image, du.last_name || ' ' || du.first_name as dest_name, du.profile_image as dest_profile_image").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_users as su ON (su.id = i.src_user_id)").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".tbd_users as du ON (du.id = i.dest_user_id)").Where("i.src_user_id = ? OR i.dest_user_id = ?", userId, userId).Order("created_date DESC").Limit(50).Scan(&invoices)

	for i := range invoices {
		invoices[i].ListType = "SEND"
		if invoices[i].DestUserId == userId {
			invoices[i].ListType = "RECEIVE"
			invoices[i].DestProfileImage = ""
		} else {
			invoices[i].SrcProfileImage = ""
		}
	}
	return core.Resolve(200, c, core.Response("success", invoices))
}

func (u *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	req := new(account_model.ReqInvoiceCreate)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	userId := oauth.GetSessionUserId(c)

	if userId == req.DestUserId {
		return core.Resolve(400, c, core.Response("dest id is invalid"))
	}

	if req.Amount < 10 {
		return core.Resolve(400, c, core.Response("invalid amount"))
	}

	db := database.DBconn

	if req.Type == "" {
		req.Type = "USER_TO_USER"
	}

	currentUser := user.ById(userId)
	if currentUser == nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	srcAccount, err := account_model.GetDefaultAccount(userId)
	if err != nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	destAccount, err := account_model.GetDefaultAccount(req.DestUserId)
	if err != nil {
		return core.Resolve(400, c, core.Response("wrong information"))
	}

	var inv = account_model.Invoice{
		SrcUserId:     userId,
		SrcAccountNo:  srcAccount.AccountNo,
		DestUserId:    req.DestUserId,
		DestAccountNo: destAccount.AccountNo,
		Amount:        req.Amount,
		Type:          req.Type,
		Status:        "NEW",
		Description:   req.Description,
		BpInvoiceBody: req.BpInvoiceBody,
		DueDate:       data.LocalTime(time.Now().AddDate(0, 0, 1)),
	}

	if err := db.Create(&inv).Error; err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	// notif_text := currentUser.FirstName + "-с нэхэмжлэх ирлээ"
	// go client.SendNotif(req.DestUserId, notif_text, notif_text, "invoice", string(rune(inv.Id)))

	return core.Resolve(200, c, core.Response("success"))
}

func (u *InvoiceHandler) PayInvoice(c *fiber.Ctx) error {
	req := new(account_model.ReqInvoicePay)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	destUserId := oauth.GetSessionUserId(c)

	db := database.DBconn

	var invoice = account_model.Invoice{}

	err := db.Where("id = ? AND dest_user_id = ? AND status = 'NEW'", req.Id, destUserId).First(&invoice).Error
	if err != nil {
		return core.Resolve(400, c, core.Response("invoice not found"))
	}

	destAccount, err := account_model.GetAccount(invoice.DestAccountNo)
	if err != nil {
		return core.Resolve(400, c, core.Response("account not found"))
	}

	if destAccount.Balance < invoice.Amount {
		return core.Resolve(400, c, core.Response("balance not available"))
	}

	journal_no := uuid.New().String()
	tdbMirrorAccount := "2000000003"

	if err := MakeTransaction(invoice.Amount, journal_no, invoice.DestAccountNo, tdbMirrorAccount, "Нэхэмжлэл төлөлт", "", "", "", "", ""); err != nil {
		return core.Resolve(400, c, core.Response("transaction error1:"+err.Error()))
	}

	if err := MakeTransaction(invoice.Amount, journal_no, tdbMirrorAccount, invoice.SrcAccountNo, "Нэхэмжлэл төлөлт", "", "", "", "", ""); err != nil {
		return core.Resolve(400, c, core.Response("transaction error2:"+err.Error()))
	}

	invoice.Status = "PAID"
	db.Save(&invoice)

	return core.Resolve(200, c, core.Response("success"))
}

func (u *InvoiceHandler) CancelInvoice(c *fiber.Ctx) error {
	req := new(account_model.ReqInvoiceCancel)
	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	destUserId := oauth.GetSessionUserId(c)

	db := database.DBconn

	var invoice = account_model.Invoice{}

	err := db.Where("id = ? AND dest_user_id = ? AND status = 'NEW'", req.Id, destUserId).First(&invoice).Error
	if err != nil {
		return core.Resolve(400, c, core.Response("invoice not found"))
	}

	invoice.Status = "CANCELLED"
	if err := db.Save(&invoice).Error; err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response("success"))
}
