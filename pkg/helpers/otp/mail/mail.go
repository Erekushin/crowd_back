package mail

import (
	"os"

	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
)

type SendEmail struct {
	SenderId     uint                   `json:"sender_id"`
	To           []string               `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateId   uint                   `json:"template_id"`
	TemplateData map[string]interface{} `json:"template_data"`
}

func Send(sender_id, template_id uint, to, subject string, body map[string]interface{}) (err error) {
	email := SendEmail{}
	email.SenderId = sender_id
	email.TemplateId = template_id
	email.To = append(email.To, to)
	email.Subject = subject
	email.TemplateData = body
	data, err := convertor.InterfaceToMap(email)
	if err != nil {
		return err
	}
	go client.Request(os.Getenv("EMAIL_SERVER_URL"), "POST", data)
	return err
}
