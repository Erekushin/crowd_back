package auth

import (
	"os"
	"strings"

	"crowdfund/pkg/base/organization"
	"crowdfund/pkg/base/user"
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/otp"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

func (*AuthHandler) Login(c *fiber.Ctx) error {

	req := new(LoginReq)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if req.Username == "" || req.Password == "" {
		return core.Resolve(400, c, core.Response("username and password is required"))
	}

	u := new(user.User)
	u.Username = strings.ToLower(req.Username)
	if err := u.FindByUsername(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	encrypted := helpers.GeneratePassword(req.Password)
	if u.Password != encrypted {
		return core.Resolve(400, c, core.Response("wrong user information"))
	}

	result := make(map[string]interface{})
	result["user"] = u

	session := oauth.GenerateSession(c, u.Id)
	result["authorization"] = session

	orgs, _ := user.OrgList(int(u.Id))
	result["organizations"] = orgs

	role := organization.GetUserRoleWithActions(session.OrgId, session.UserId)
	result["role"] = role
	// term := models.Terminals{}
	// if err := db.Where("serial_no = ?", req.Serialno).Find(&term).Error; err != nil {
	// 	return ErrorDetails(err)
	// }

	// result["terminal"] = term
	return c.JSON(result)
}

func (*AuthHandler) Identify(c *fiber.Ctx) error {
	text := c.Query("text")

	if text == "" {
		return core.Resolve(400, c, core.Response("identity /email or phone_no/ text is required"))
	}

	u := new(user.User)

	if core.ValidEmail(text) {
		u.Email = text
		if u.ExistWithEmail() {
			return core.Resolve(400, c, core.Response("user email already exists"))
		}
	} else if core.ValidPhone(text) {
		u.PhoneNo = text
		if u.ExistWithPhone() {
			return core.Resolve(400, c, core.Response("user phone_no already exists"))
		}
	} else {
		return core.Resolve(400, c, core.Response("identity text is invalid"))
	}

	if err := otp.SendOtp(text); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*AuthHandler) Register(c *fiber.Ctx) error {
	data := new(RegisterReq)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := otp.CheckOtp(data.Identity, data.Otp); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	user := new(user.User)
	data.Identity = strings.ToLower(data.Identity)
	if core.ValidEmail(data.Identity) {
		user.Email = data.Identity
		user.IsConfirmedEmail = 1
		if user.ExistWithEmail() {
			return core.Resolve(400, c, core.Response("user email already exists"))
		}
	} else if core.ValidPhone(data.Identity) {
		user.PhoneNo = data.Identity
		user.IsConfirmedPhoneNo = 1
		if user.ExistWithPhone() {
			return core.Resolve(400, c, core.Response("user phone_no already exists"))
		}
	} else {
		return core.Resolve(400, c, core.Response("identity phone_no or email is invalid"))
	}

	user.Username = data.Identity
	user.Password = helpers.GeneratePassword(data.Password)
	if err := user.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	result := make(map[string]interface{})
	result["user"] = user
	result["authorization"] = oauth.GenerateSession(c, user.Id)

	return core.Resolve(200, c, core.Response("success", result))
}

func (*AuthHandler) ChangeOrg(c *fiber.Ctx) error {
	orgId := convertor.StringToInt(c.Query("org_id"))

	if orgId == 0 {
		return core.Resolve(400, c, core.Response("org_id is required"))
	}

	org := organization.FindById(uint(orgId))
	if org == nil {
		return core.Resolve(400, c, core.Response("data not found"))
	}

	if err := oauth.ChangeOrg(c, org.Id); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	userId := oauth.GetSessionUserId(c)

	role := organization.GetUserRoleWithActions(org.Id, userId)

	return core.Resolve(200, c, core.Response("success", role))
}

func (*AuthHandler) ResetPasswordOtp(c *fiber.Ctx) error {
	identity := c.Query("identity")

	if identity == "" {
		return core.Resolve(400, c, core.Response("identity /email or phone_no/ text is required"))
	}

	u := new(user.User)
	userExists := false
	if core.ValidEmail(identity) {
		u.Email = identity
		userExists = u.ExistWithEmail()
	} else if core.ValidPhone(identity) {
		u.PhoneNo = identity
		userExists = u.ExistWithPhone()
	} else {
		return core.Resolve(400, c, core.Response("identity /email or phone_no/ text is invalid"))
	}

	if !userExists {
		return core.Resolve(400, c, core.Response("wrong user identity"))
	}

	if err := otp.SendOtp(identity); err != nil {
		if err.Error() == "MESSAGE_SENT" {
			return core.Resolve(201, c, core.Response("OTP sent"))
		}
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*AuthHandler) ResetPassword(c *fiber.Ctx) error {
	otptext := convertor.StringToInt(c.Query("otp"))
	identity := c.Query("identity")
	password := c.Query("password")

	if otptext == 0 {
		return core.Resolve(400, c, core.Response("otp is required"))
	}

	if identity == "" {
		return core.Resolve(400, c, core.Response("identity /email or phone_no/ is required"))
	}

	if password == "" {
		return core.Resolve(400, c, core.Response("password is required"))
	}

	u := new(user.User)
	userExists := false
	if core.ValidEmail(identity) {
		u.Email = identity
		userExists = u.ExistWithEmail()
	} else if core.ValidPhone(identity) {
		u.PhoneNo = identity
		userExists = u.ExistWithPhone()
	} else {
		return core.Resolve(400, c, core.Response("identity /email or phone_no/ text is invalid"))
	}

	if !userExists {
		return core.Resolve(400, c, core.Response("wrong user identity"))
	}

	if err := otp.CheckOtp(identity, uint(otptext)); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	if err := u.UpdatePassword(password); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	// return core.Resolve(200, c, core.Response())
	return core.Resolve(200, c, core.Response("success"))
}

func (*AuthHandler) QrLogin(c *fiber.Ctx) error {
	userId := oauth.GetSessionUserId(c)
	serialNo := c.Query("serial_no")

	if serialNo == "" {
		return core.Resolve(400, c, core.Response("serial_no is required"))
	}

	u := user.User{Id: userId}
	if err := u.FindById(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	authInfo := make(map[string]interface{})
	authInfo["user"] = u

	session := oauth.GenerateSession(c, u.Id)
	authInfo["authorization"] = session

	orgs, _ := user.OrgList(int(u.Id))
	authInfo["organizations"] = orgs

	role := organization.GetUserRoleWithActions(session.OrgId, session.UserId)
	authInfo["role"] = role

	payload := make(map[string]interface{})
	payload["to"] = serialNo
	payload["event"] = "goQrLogin"
	payload["data"] = authInfo

	go client.Request(os.Getenv("SOCKET_URL"), "POST", payload, nil)

	return core.Resolve(200, c, core.Response("success"))
}
