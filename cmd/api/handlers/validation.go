package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/tshivanshu9/budget-be/common"
)

func (h *Handler) ValidateBodyRequest(c *echo.Context, payload interface{}) []*common.ValidationError {
	var errors []*common.ValidationError
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(payload)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			reflected := reflect.ValueOf(payload)
			for _, validationError := range validationErrors {
				field, _ := reflected.Type().FieldByName(validationError.StructField())
				key := field.Tag.Get("json")
				if key == "" {
					key = strings.ToLower(validationError.StructField())
				}
				condition := validationError.Tag()
				keyToTitleCase := strings.Replace(key, "_", " ", -1)
				param := validationError.Param()
				errMessage := keyToTitleCase + " field is " + condition

				switch condition {
				case "required":
					errMessage = keyToTitleCase + " is required"
				case "email":
					errMessage = keyToTitleCase + " must be a valid email"
				case "min":
					errMessage = fmt.Sprintf("%s must be atleast %s", keyToTitleCase, param)
				}

				fmt.Println(validationError.Field())
				currentValidationError := common.ValidationError{
					Error:     errMessage,
					Key:       key,
					Condition: condition,
				}
				errors = append(errors, &currentValidationError)
			}
		}
	}

	return errors
}
