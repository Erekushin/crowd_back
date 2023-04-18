package sms

import (
	"os"

	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
)

type SendSms struct {
	PhoneNumbers []string `json:"phone_numbers"`
	MessageValue string   `json:"message_value"`
	SmsType      uint     `json:"sms_type"`
}

func Send(phone, text string) (err error) {
	sms := SendSms{}
	sms.PhoneNumbers = append(sms.PhoneNumbers, phone)
	sms.MessageValue = text
	sms.SmsType = 10
	data, err := convertor.InterfaceToMap(sms)
	header := make(map[string]string)
	header["message_code"] = os.Getenv("OTP_SMS_MESSAGE_CODE")
	go client.Request(os.Getenv("SMS_SERVER_URL"), "POST", data, header)
	if err != nil {
		return err
	}
	return err
}
