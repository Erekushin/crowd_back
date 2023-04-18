package oauth

import (
	"errors"

	"crowdfund/pkg/database"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/gofiber/fiber/v2"
)

var srv *server.Server

func Init() error {
	manager := manage.NewDefaultManager()

	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()

	var clients []Client

	db := database.DBconn

	result := db.Find(&clients)

	if result.Error != nil {
		return result.Error
	}

	for _, each := range clients {
		clientStore.Set(each.ClientId, &models.Client{
			ID:     each.ClientId,
			Secret: each.ClientSecret,
			Domain: "",
		})
	}

	manager.MapClientStorage(clientStore)

	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	return nil
}

func GenerateSession(c *fiber.Ctx, userId uint) *Session {
	db := database.DBconn
	session := Session{}

	if err := db.Find(&session, "user_id = ?", userId).Error; err != nil {
		return nil
	}

	if session.Id != 0 {
		session.Extend()
		return &session
	}

	return CreateSession(userId)
}

func ChangeOrg(c *fiber.Ctx, orgId uint) error {
	userId := GetSessionUserId(c)

	db := database.DBconn
	session := Session{}

	if err := db.Find(&session, "user_id = ?", userId).Error; err != nil {
		return err
	}

	if session.Id == 0 {
		return errors.New("session not found")
	}

	session.OrgId = orgId
	session.Save()
	return nil
}
