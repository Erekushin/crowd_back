package core

import (
	"net/mail"
	"regexp"

	"crowdfund/pkg/helpers"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func Validate(data interface{}) []*ValidationError {

	var errors []*ValidationError
	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var ve ValidationError

			ve.Field = helpers.ToSnakeCase(err.StructField())
			ve.Tag = err.Tag()
			ve.Param = err.Param()
			errors = append(errors, &ve)
		}
	}

	return errors
}

func ValidEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	return err == nil
}

func ValidPhone(str string) bool {
	ok, _ := regexp.MatchString(`^\d{8}$`, str)
	return ok
}
