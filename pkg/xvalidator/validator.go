package xvalidator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(input interface{}) error {
	if err := validate.Struct(input); err != nil {
		return formatValidationError(err)
	}

	return nil
}

func formatValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	if len(validationErrors) == 0 {
		return err
	}

	fieldError := validationErrors[0]
	field := strings.ToLower(fieldError.Field())

	switch fieldError.Tag() {
	case "required":
		return fmt.Errorf("%s is required", field)
	case "email":
		return fmt.Errorf("%s must be a valid email", field)
	case "min":
		return fmt.Errorf("%s must be at least %s characters", field, fieldError.Param())
	case "oneof":
		return fmt.Errorf("%s must be one of: %s", field, fieldError.Param())
	case "datetime":
		return fmt.Errorf("%s must use format YYYY-MM-DD", field)
	default:
		return fmt.Errorf("%s is invalid", field)
	}
}
