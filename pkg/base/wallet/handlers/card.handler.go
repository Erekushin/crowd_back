package account_handler

import (
	"encoding/json"
	"fmt"
	"os"

	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type CardService struct{}

func (u *CardService) CreateCardInvoice(c *fiber.Ctx) error {

	payload := map[string]interface{}{
		"user_id":      oauth.GetSessionUserId(c),
		"app_id":       convertor.StringToInt(os.Getenv("APP_ID")),
		"redirect_url": c.Query("redirect_url"),
	}
	cr, _ := json.Marshal(payload)
	result, _ := client.SendRequest(os.Getenv("CARD_URL")+"/invoice", "", "", "POST", "", []byte(cr))

	var res = map[string]interface{}{}
	err := json.Unmarshal(result, &res)
	if err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response(res["message"], res["result"]))
}

func (u *CardService) GetCardList(c *fiber.Ctx) error {
	payload := map[string]interface{}{
		"user_id": oauth.GetSessionUserId(c),
		"app_id":  convertor.StringToInt(os.Getenv("APP_ID")),
	}
	cr, _ := json.Marshal(payload)
	result, _ := client.SendRequest(os.Getenv("CARD_URL")+"/list", "", "", "POST", "", []byte(cr))

	var res = map[string]interface{}{}
	err := json.Unmarshal(result, &res)
	if err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response(res["message"], res["result"]))
}

func (u *CardService) CardDelete(c *fiber.Ctx) error {

	req := new(account_model.ReqCardDelete)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	payload := map[string]interface{}{
		"user_id": oauth.GetSessionUserId(c),
		"app_id":  convertor.StringToInt(os.Getenv("APP_ID")),
		"id":      req.Id,
	}

	cr, _ := json.Marshal(payload)
	result, _ := client.SendRequest(os.Getenv("CARD_URL")+"/delete", "", "", "POST", "", []byte(cr))

	var res = map[string]interface{}{}
	err := json.Unmarshal(result, &res)
	if err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response(res["message"], res["result"]))
}

func (u *CardService) CardDeposit(c *fiber.Ctx) error {

	req := new(account_model.ReqCardPay)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if req.Amount < 1 {
		return core.Resolve(400, c, core.Response("invalid amount"))
	}

	if req.ChargeAmount < 1 {
		return core.Resolve(400, c, core.Response("invalid charge amount"))
	}

	fmt.Println("CardDeposit req:", req)

	userId := oauth.GetSessionUserId(c)
	payload := map[string]interface{}{
		"app_id":         convertor.StringToInt(os.Getenv("APP_ID")),
		"user_id":        userId,
		"type":           "charge",
		"hash":           req.Hash,
		"card_token_id":  req.CardTokenId,
		"device_type":    req.DeviceType,
		"amount":         req.Amount,
		"charge_percent": req.ChargePercent,
		"charge_amount":  req.ChargeAmount,
	}

	fmt.Println("CardDeposit payload:", payload)

	cr, _ := json.Marshal(payload)
	result, _ := client.SendRequest(os.Getenv("CARD_URL")+"/pay", "", "", "POST", "", []byte(cr))

	response := new(account_model.ResCardDepsit)
	err := json.Unmarshal(result, response)
	if err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}
	fmt.Println("CardDeposit response:", response)

	if response.Code != 200 {
		return core.Resolve(500, c, core.Response(response.Message))
	}
	fmt.Println("response code:", response.Code)

	fmt.Println("ref_no:", response.Result.RefNo)

	err = depo_account(userId, req.ChargeAmount, response.Result.RefNo, "Deposit: "+req.CardNo)
	if err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response("success"))
}
