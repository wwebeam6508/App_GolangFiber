package common

import (
	"PBD_backend_go/exception"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func Validate(modelValidate interface{}) error {
	validate := validator.New()
	err := validate.Struct(modelValidate)
	if err != nil {
		var messages []map[string]interface{}
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, map[string]interface{}{
				"field":   err.Field(),
				"message": "this field is " + err.Tag(),
			})
		}

		jsonMessage, _ := json.Marshal(messages)

		return exception.ValidationError{
			Message: string(jsonMessage),
		}
	}
	return nil
}
