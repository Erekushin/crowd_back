package core

import (
	"encoding/json"
	"os"

	"crowdfund/pkg/base/message"
	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/data"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type Log struct {
	Id          uint           `json:"id,string" gorm:"primaryKey"`
	RequestId   string         `json:"request_id" gorm:"type:varchar(100)"`
	MessageCode uint           `json:"message_code"`
	Path        string         `json:"path" gorm:"type:varchar(200)"`
	Method      string         `json:"method" gorm:"type:varchar(20);"`
	Payload     datatypes.JSON `json:"payload"`
	Type        string         `json:"type" gorm:"type:varchar(20);default:REQUEST"`
	StatusCode  uint           `json:"status_code" gorm:"type:int4;default:200"`
	SessionId   uint           `json:"session_id"`
	CreatedAt   data.LocalTime `json:"created_date" gorm:"autoCreateTime"`
	CreatedBy   uint           `json:"created_by,omitempty"`
}

func (*Log) TableName() string {
	return os.Getenv("DB_SCHEMA") + ".tbs_logs"
}

func GetMessage(c *fiber.Ctx) *message.Message {
	cs := c.Locals("message")
	msg, _ := cs.(*message.Message)
	return msg
}

func SaveRequestLog(c *fiber.Ctx) {
	requestId, _ := convertor.InterfaceToString(c.Locals("requestid"), "requestid is null")
	session := oauth.GetSession(c)
	message := GetMessage(c)

	body := make(map[string]interface{})

	if message.IsQuery == 1 {
		c.QueryParser(&body)
	} else {
		json.Unmarshal([]byte(c.Body()), &body)
	}

	jsonBody, _ := json.Marshal(body)

	log := Log{
		RequestId:   requestId,
		MessageCode: message.Code,
		Path:        message.Path,
		Method:      message.Method,
		Payload:     jsonBody,
		SessionId:   session.Id,
		CreatedBy:   session.UserId,
	}
	if message.IsDbLog == 1 {
		db := database.DBconn
		db.Create(&log)
	}

	if message.IsFileLog == 1 {
		// implement write log to file
	}
}

func SaveResponseLog(statusCode uint, c *fiber.Ctx, res *ApiResponse) {
	requestId, _ := convertor.InterfaceToString(c.Locals("requestid"), "requestid is null")
	session := oauth.GetSession(c)
	message := GetMessage(c)

	jsonBody, _ := json.Marshal(res)
	log := Log{
		RequestId:   requestId,
		MessageCode: message.Code,
		Path:        message.Path,
		Method:      message.Method,
		Payload:     jsonBody,
		SessionId:   session.Id,
		CreatedBy:   session.UserId,
		Type:        "RESPONSE",
		StatusCode:  statusCode,
	}
	if message.IsDbLog == 1 {
		db := database.DBconn
		db.Create(&log)
	}

	if message.IsFileLog == 1 {
		// implement write log to file
	}
}
