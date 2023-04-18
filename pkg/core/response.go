package core

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

func Response(responseDetails ...interface{}) *ApiResponse {
	var (
		result  interface{}
		message string
	)
	for index, val := range responseDetails {
		switch index {
		case 0:
			message = val.(string)

		case 1:
			result = val
		}
	}

	if message == "" {
		message = "success"
	}

	if result == nil {
		result = fiber.Map{}
	}

	return &ApiResponse{
		Message: message,
		Result:  result,
	}
}
