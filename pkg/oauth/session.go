package oauth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Session struct {
	Id             uint           `json:"-" gorm:"primaryKey"`
	UserId         uint           `json:"user_id,string,omitempty"`
	OrgId          uint           `json:"org_id,string,omitempty"`
	TerminalId     uint           `json:"terminal_id,string,omitempty"`
	Token          string         `json:"token" gorm:"index"`
	Expires        time.Time      `json:"-"`
	ExpireAsString string         `json:"expires" gorm:"-"`
	ExpiresIn      uint           `json:"expires_in" gorm:"-"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (*Session) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".oauth_sessions2"
}

func (s *Session) Save() error {
	db := database.DBconn
	return db.Save(s).Error
}

func (s *Session) Extend() {
	s.Expires = time.Now().Add(time.Hour * time.Duration(convertor.StringToInt(os.Getenv("OAUTH_TOKEN_EXP_HOURS"))))
	s.ExpireAsString = convertor.TimeToDateTimeString(s.Expires)
	s.ExpiresIn = uint(s.Expires.Sub(time.Now()).Seconds())
	s.Save()
}

func GetSession(c *fiber.Ctx) (s *Session) {
	cs := c.Locals("session")
	session, ok := cs.(*Session)
	if !ok {
		return &Session{}
	}
	return session
}

func GetSessionUserId(c *fiber.Ctx) uint {
	session := GetSession(c)
	return session.UserId
}

func GetSessionOrgId(c *fiber.Ctx) uint {
	session := GetSession(c)
	return session.OrgId
}

func GetSessionTerminalId(c *fiber.Ctx) uint {
	session := GetSession(c)
	return session.TerminalId
}

func CreateSession(userId uint) *Session {
	ctx := context.Background()
	gt := oauth2.GrantType("client_credentials")

	tgr := &oauth2.TokenGenerateRequest{
		ClientID:       os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret:   os.Getenv("OAUTH_CLIENT_SECRET"),
		Request:        &http.Request{},
		AccessTokenExp: time.Hour * time.Duration(convertor.StringToInt(os.Getenv("OAUTH_TOKEN_EXP_HOURS"))),
	}

	ti, _ := srv.GetAccessToken(ctx, gt, tgr)

	session := Session{
		UserId:    userId,
		Token:     strings.ToLower(ti.GetAccess()),
		Expires:   time.Now().Add(ti.GetAccessExpiresIn()),
		ExpiresIn: uint(ti.GetAccessExpiresIn() / 1e9),
	}
	session.ExpireAsString = convertor.TimeToDateTimeString(session.Expires)
	session.Save()

	return &session
}
