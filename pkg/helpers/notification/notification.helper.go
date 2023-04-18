package notification

import (
	"os"

	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
)

type Notifications struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	Tokens   []string `json:"tokens"`
	Platform int      `json:"platform"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	ApiKey   string   `json:"api_key"`
}

func Send(device_token, title, msg string) (err error) {

	var reqs Notifications
	var req Notification

	req.Tokens = append(req.Tokens, device_token)
	req.Title = title
	req.Message = msg
	req.Platform = 2
	req.ApiKey = os.Getenv("NOTIF_APKEY")

	reqs.Notifications = append(reqs.Notifications, req)

	data, err := convertor.InterfaceToMap(reqs)
	client.Request(os.Getenv("NOTIF_SERVER_URL"), "POST", data)
	return err
}
