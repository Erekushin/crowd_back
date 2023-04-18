package crowdfund

import (
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"crowdfund/pkg/base/organization"
	"crowdfund/pkg/base/user"
	account_model "crowdfund/pkg/base/wallet/models"

	"gorm.io/gorm"
)

func CreateCrowndfund(req *ReqCrowdfundCreate, ctx *fiber.Ctx) error {

	db := database.DBconn
	s := oauth.GetSession(ctx)
	crowdfund := new(Crowdfund)
	crowdfund.OrgId = s.OrgId
	crowdfund.OrgName = organization.GetName(s.OrgId)
	crowdfund.CreatedBy = s.UserId
	crowdfund.Name = req.Name
	crowdfund.Description = req.Description
	crowdfund.ImgBase64 = req.ImgBase64
	crowdfund.CategoryId = req.CategoryId
	crowdfund.Amount = req.Amount
	crowdfund.ProfitPrecent = req.ProfitPrecent
	crowdfund.Status = "PENDING"
	crowdfund.StartDate = convertor.DateStringToTime(req.StartDate)
	crowdfund.EndDate = convertor.DateStringToTime(req.EndDate)
	crowdfund.IntroductionText = req.IntroductionText
	crowdfund.RiskText = req.RiskText
	return db.Create(&crowdfund).Error
}

func (c *Crowdfund) Update() error {
	db := database.DBconn
	return db.Updates(&c).Error
}

func (c *Crowdfund) Delete() error {
	db := database.DBconn
	return db.Delete(c).Error
}

func List(c *fiber.Ctx) (*data.Pagination, error) {
	var totalRow int64
	db := database.DBconn
	users := make([]Crowdfund, 0)

	orgId := oauth.GetSessionOrgId(c)
	org := organization.FindById(orgId)
	if org == nil {
		return nil, errors.New("organization not found")
	}

	tx := db.Model(Crowdfund{})
	if org.TypeId == 2 {
		tx.Where("org_id=?", orgId)
	}
	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Order("created_at desc").Offset(p.Offset).Limit(p.PageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	p.Items = users
	return p, nil
}

func ListConfirmed(c *fiber.Ctx) (*data.Pagination, error) {
	SearchText := strings.ToLower(c.Query("search_text"))
	CategoryId := c.Query("category_id")

	var fields []string
	var values []any

	var totalRow int64
	db := database.DBconn
	crowdfund := make([]Crowdfund, 0)

	if SearchText != "" {
		fields = append(fields, "lower(c.name)  LIKE ?")
		values = append(values, "%"+SearchText+"%")
	}

	if CategoryId != "" {
		fields = append(fields, "c.category_id = ?")
		values = append(values, convertor.StringToInt(CategoryId))
	}
	currentDate := time.Now()
	fields = append(fields, "c.status = ?")
	values = append(values, "CONFIRMED")
	fields = append(fields, "c.end_date >= ?")
	values = append(values, currentDate)

	tx := db.Table(os.Getenv("DB_SCHEMA")+".cwd_crowdfund as c").Select("c.*, a.balance, count(distinct(cu.user_id)) user_cnt").Joins("INNER JOIN "+os.Getenv("DB_SCHEMA")+".tbw_accounts as a ON (a.owner_id = c.id)").Joins("LEFT JOIN "+os.Getenv("DB_SCHEMA")+".cwd_crowdfund_user as cu ON (cu.crowdfund_id = c.id)").Where(strings.Join(fields, " AND "), values...).Order("c.created_at desc").Group("c.id, a.balance")

	tx.Count(&totalRow)

	p := data.Paginate(c, totalRow)

	err := tx.Offset(p.Offset).Limit(p.PageSize).Find(&crowdfund).Error
	if err != nil {
		return nil, err
	}
	p.Items = crowdfund
	return p, nil
}

func ConfirmCrowndfund(req *ReqCrowdfundId, ctx *fiber.Ctx) error {
	db := database.DBconn
	fund := new(Crowdfund)

	if err := db.First(&fund, "id=? AND status='PENDING'", req.Id).Error; err != nil {
		return errors.New(" Төслийн мэдээлэл олдсонгүй")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var account account_model.Account
		account.OwnerId = fund.Id
		account.Name = fund.Name
		account.Label = fund.Name
		account.TypeId = 30
		account.TypeName = "Төсөл"

		if err := account.Create(); err != nil {
			return errors.New("Төслийн данс үүсгэж чадсангүй:" + err.Error())
		}

		if err := db.Model(&Crowdfund{}).Where("id = ?", req.Id).Update("status", "CONFIRMED").Error; err != nil {
			return errors.New("Төсөл үүсгэж чадсангүй:" + err.Error())
		}
		return nil
	})

}

func CancelCrowdfund(req *ReqCrowdfundId, ctx *fiber.Ctx) error {
	db := database.DBconn
	fund := new(Crowdfund)
	if err := db.First(&fund, "id=? AND status='PENDING'", req.Id).Error; err != nil {
		return errors.New(" Төслийн мэдээлэл олдсонгүй")
	}

	return db.Model(&Crowdfund{}).Where("id = ?", req.Id).Update("status", "CANCELED").Error
}

func Info(req *ReqCrowdfundId, c *fiber.Ctx) (any, error) {
	var err error
	db := database.DBconn
	fund := Crowdfund{}

	if err := db.First(&fund, "id=? AND status='CONFIRMED'", req.Id).Error; err != nil {
		return err, err
	}
	u := new(user.User)
	u.Id = fund.CreatedBy
	if err := u.FindById(); err != nil {
		return err, err
	}

	result := make(map[string]interface{})

	result["crowdfund"] = fund
	result["user"] = u

	orgs, _ := user.OrgListById(int(u.Id), fund.OrgId)
	result["organizations"] = orgs

	return result, err
}
