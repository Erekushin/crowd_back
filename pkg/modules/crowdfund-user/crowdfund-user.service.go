package crowdfunduser

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"
	"errors"
	"os"

	account "crowdfund/pkg/base/wallet/handlers"
	account_model "crowdfund/pkg/base/wallet/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCrowdfundUser(req *ReqCrowdfundUserCreate, ctx *fiber.Ctx) error {
	db := database.DBconn
	return db.Transaction(func(tx *gorm.DB) error {
		srcAccount, err := account_model.GetDefaultAccount(req.UserId)
		if err != nil {
			return core.Resolve(400, nil, core.Response("invalid src account"))
		}

		if srcAccount.Balance < float32(req.Amount) {
			return errors.New("Дансны үлдэгдэл хүрэлцэхгүй байна.")
		}

		destAccount, err := account_model.GetDefaultAccount(req.CrowdfundId)
		if err != nil {
			return core.Resolve(400, nil, core.Response("invalid dest account"))
		}

		tdbMirrorAccount := "2000000003"
		journal_no := uuid.New().String()

		if err := account.MakeTransaction(float32(req.Amount), journal_no, srcAccount.AccountNo, tdbMirrorAccount, "", "", "", "", "", ""); err != nil {
			return core.Resolve(400, nil, core.Response("transaction error1:"+err.Error()))
		}

		if err := account.MakeTransaction(float32(req.Amount), journal_no, tdbMirrorAccount, destAccount.AccountNo, "", "", "", "", "", ""); err != nil {
			return core.Resolve(400, nil, core.Response("transaction error2:"+err.Error()))
		}

		s := oauth.GetSession(ctx)
		fundUser := new(CrowdfundUser)
		fundUser.CreatedBy = s.UserId
		fundUser.CrowdfundId = req.CrowdfundId
		fundUser.UserId = req.UserId
		fundUser.Amount = req.Amount
		return db.Create(&fundUser).Error
	})

}

func CrowdfundUserList(CrowdfundID int, c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	users := make([]CrowdfundUsers, 0)

	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_users as u").Select("u.id as user_id, u.first_name, u.last_name, u.reg_no, c.id crowdfund_id, sum(cu.amount) amount").Joins("INNER JOIN "+os.Getenv("DB_SCHEMA")+".cwd_crowdfund_user as cu ON (cu.user_id = u.id)").Joins("INNER JOIN "+os.Getenv("DB_SCHEMA")+".cwd_crowdfund as c ON (cu.crowdfund_id = c.id)").Where("c.id=?", CrowdfundID).Group("u.id, first_name, last_name, reg_no, c.id")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	p.Items = users
	return p, nil
}

func UserCrowdfundList(UserId uint, c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	crowdfund := make([]UserCrowdfund, 0)

	tx := db.Table(os.Getenv("DB_SCHEMA")+".tbd_users as u").Select("c.*, sum(cu.amount) invested_amount").Joins("INNER JOIN "+os.Getenv("DB_SCHEMA")+".cwd_crowdfund_user as cu ON (cu.user_id = u.id)").Joins("INNER JOIN "+os.Getenv("DB_SCHEMA")+".cwd_crowdfund as c ON (cu.crowdfund_id = c.id)").Where("u.id=?", UserId).Group("c.id")
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&crowdfund).Error
	if err != nil {
		return nil, err
	}
	p.Items = crowdfund
	return p, nil
}
