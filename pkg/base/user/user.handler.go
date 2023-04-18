package user

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"crowdfund/pkg/core"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/otp"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func (*UserHandler) Find(c *fiber.Ctx) error {
	searchText, _ := url.QueryUnescape(c.Query("search_text"))
	if searchText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "search text is null", "result": fiber.Map{}})
	}

	user := User{}
	db := database.DBconn

	if len(searchText) == 8 {
		db.Unscoped().First(&user, "id=?", searchText)
	}

	if len(searchText) == 13 {
		db.Unscoped().First(&user, "civil_id=?", searchText)
	}

	if ok, _ := regexp.MatchString("^[0-9]{12}$", searchText); ok {
		db.Unscoped().First(&user, "civil_id=?", searchText)
	}

	if ok, _ := regexp.MatchString("^[а-яА-ЯөӨүҮ]{2}[0-9]{8}$", searchText); ok {
		db.Unscoped().First(&user, "reg_no=?", strings.ToLower(searchText))
	}

	if user.Id == 0 {
		body := make(map[string]interface{})
		body["search_text"] = searchText
		response := client.Request(os.Getenv("URL_USER"), "POST", body)
		if response.Code != 200 {
			return c.JSON(response)
		}

		convertor.MapToStruct(response.Result, &user)
		user.Username = user.RegNo
		user.Password = helpers.GeneratePassword("81dc9bdb52d04dc20036dbd8313ed055") // 1234
		user.CreatedBy = oauth.GetSessionUserId(c)
		if err := db.Create(&user).Error; err != nil {
			return core.Resolve(500, c, core.Response(err.Error()))
		}
	}
	return core.Resolve(200, c, core.Response("success", user))
}

func (*UserHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = List(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*UserHandler) Update(c *fiber.Ctx) error {
	data := new(User)
	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.UpdatedBy = oauth.GetSessionUserId(c)
	if err := data.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response())
}

func (*UserHandler) Delete(c *fiber.Ctx) error {
	data := new(User)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.DeletedBy = oauth.GetSessionUserId(c)
	if err := data.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*UserHandler) ChangePassword(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	user := User{Id: userId}
	oldPassword := c.Query("old")
	newPassword := c.Query("new")
	if oldPassword == "" || newPassword == "" {
		return core.Resolve(400, c, core.Response("old and new passwords required"))
	}

	err := user.ChangePassword(oldPassword, newPassword)
	if err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*UserHandler) ChangeUsername(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	user := User{Id: userId}
	username := strings.ToLower(c.Query("text"))
	if username == "" {
		return core.Resolve(400, c, core.Response("username is required"))
	}
	err := user.ChangeUsername(username)
	if err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*UserHandler) FindPhone(c *fiber.Ctx) error {
	searchText, _ := url.QueryUnescape(c.Query("search_text"))
	if searchText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "search text is null", "result": fiber.Map{}})
	}

	if ok, _ := regexp.MatchString("^[0-9]{8}$", searchText); !ok {
		return core.Resolve(404, c, core.Response("phone is not valid"))
	}

	user := ApiUser{}
	db := database.DBconn
	db.Model(User{}).Take(&user, "phone_no=?", searchText)

	if user.Id == 0 {
		return core.Resolve(404, c, core.Response("user not found"))
	}

	user.FirstName = helpers.Mask(user.FirstName, 1)
	user.LastName = helpers.Mask(user.LastName, 1)
	return core.Resolve(200, c, core.Response("success", user))
}

func (*UserHandler) ChangeEmailOrPhone(c *fiber.Ctx) error {

	req := new(ReqEmailOrChange)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := otp.CheckOtp(req.Identity, uint(convertor.StringToInt(req.Otp))); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	user := new(User)
	user.Id = oauth.GetSessionUserId(c)
	req.Identity = strings.ToLower(req.Identity)
	if core.ValidEmail(req.Identity) {
		user.Email = req.Identity
		user.IsConfirmedEmail = 1
		if user.ExistWithEmail() {
			return core.Resolve(400, c, core.Response("Имейл аль хэдийн бүртгэгдсэн байна"))
		}
	} else if core.ValidPhone(req.Identity) {
		user.PhoneNo = req.Identity
		user.IsConfirmedPhoneNo = 1
		if user.ExistWithPhone() {
			return core.Resolve(400, c, core.Response("Утасны дугаар аль хэдийн бүртгэгдсэн байна"))
		}
	} else {
		return core.Resolve(400, c, core.Response("Имейл/Утасны дугаар буруу байна"))
	}

	if err := user.Save(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response("success"))
}
