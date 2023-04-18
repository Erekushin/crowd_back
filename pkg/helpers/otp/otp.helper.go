package otp

import (
	"errors"
	"fmt"
	"os"
	"time"

	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/helpers/otp/mail"
	"crowdfund/pkg/helpers/otp/sms"
)

func SendOtp(identity string) (err error) {
	otp := OtpCode{}
	waitSecond := convertor.StringToInt(os.Getenv("MESSAGE_WAIT_SECONDS"))
	maxAllowedTime := time.Now().Add(-time.Second * time.Duration(waitSecond))
	cnt := otp.CheckOtp(identity, maxAllowedTime)

	if cnt > 0 {
		return fmt.Errorf("MESSAGE_WAIT_SECOND_ERROR")
	}

	otp.Code = uint(helpers.GenerateOtp())
	otp.Identity = identity

	if err = otp.Save(); err != nil {
		return err
	}

	if core.ValidEmail(identity) {
		sender_id := convertor.StringToInt(os.Getenv("OTP_EMAIL_SENDER_ID"))
		template_id := convertor.StringToInt(os.Getenv("OTP_EMAIL_TEMPLATE_ID"))
		err = mail.Send(uint(sender_id), uint(template_id), identity, "OTP CODE", map[string]interface{}{"code": otp.Code})
		if err != nil {
			return err
		}
	} else if core.ValidPhone(identity) {
		smsReq := `OTP CODE: ` + convertor.UintToString(otp.Code)
		err = sms.Send(otp.Identity, smsReq)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid identity")
	}
	return err
}

func CheckOtp(identity string, code uint) (err error) {
	otp := OtpCode{}
	if err = otp.GetLastOtp(identity); err != nil {
		return err
	}

	if otp.Code != code {
		return errors.New("wrong otp code")
	}

	return err
}
