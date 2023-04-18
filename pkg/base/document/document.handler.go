package document

import (
	"os"
	"regexp"
	"strings"

	"crowdfund/pkg/base/user"
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct{}

func (*DocumentHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	userId := oauth.GetSessionUserId(c)

	if res, err = DocumentList(c, userId); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*DocumentHandler) Update(c *fiber.Ctx) error {
	data := new(Document)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.UserId = oauth.GetSessionUserId(c)
	if err := data.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*DocumentHandler) Delete(c *fiber.Ctx) error {

	id := convertor.StringToInt(c.Query("id"))

	if id == 0 {
		return core.Resolve(400, c, core.Response("id is required"))
	}

	data := Document{Id: uint(id)}
	data.DeletedBy = oauth.GetSessionUserId(c)
	if err := data.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response())
}

func (*DocumentHandler) SetUser(c *fiber.Ctx) error {
	req := new(SetUserReq)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	err := SetUser(req.UserId, req.DocumentId)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	header := make(map[string]string)
	header["message_code"] = os.Getenv("MESSAGE_CODE_DOCUMENT_SET_USER")
	data, _ := convertor.InterfaceToMap(req)
	go client.Request(os.Getenv("URL_USER_GO"), "POST", data, header)
	return core.Resolve(200, c, core.Response())
}

func (*DocumentHandler) Find(c *fiber.Ctx) error {
	req := new(ApiFindDocument)
	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if req.CountryCode == "" {
		return core.Resolve(400, c, core.Response("country_code is required"))
	}

	document := new(Document)
	var docFindErr error
	req.CountryCode = strings.ToLower(req.CountryCode)

	if req.CountryCode == "mng" {
		if ok, _ := regexp.MatchString("^[а-яА-ЯөӨүҮ]{2}[0-9]{8}$", req.RegNo); !ok {
			return core.Resolve(400, c, core.Response("reg_no is not valid"))
		}
		if req.DocumentNumber == "" {
			document.DocumentNumber = req.RegNo
		} else {
			document.DocumentNumber = req.DocumentNumber
		}
		docFindErr = document.FindByDocumentNumber()
	} else {
		if req.DocumentNumber == "" {
			return core.Resolve(400, c, core.Response("document_number is required"))
		}
		document.Hash = helpers.GenerateHashDocument(req.FirstName, req.LastName, req.BirthDate, req.Gender, req.CountryCode, req.DocumentNumber)
		docFindErr = document.FindByHash()
	}

	session := oauth.GetSession(c)
	if docFindErr != nil {
		return core.Resolve(200, c, core.Response("success", FindDocumentFromCore(session, req)))
	}

	u := user.User{Id: document.UserId}
	u.FindById()
	u.MergeSessionUser(session)

	result := map[string]interface{}{
		"document": document,
		"user":     u,
	}
	return core.Resolve(200, c, core.Response("success", result))

}
