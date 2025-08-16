package validators

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator structure
type Validator struct {
	*validator.Validate
}

// NewValidator creates and configures a new validator
func NewValidator() *Validator {
	v := validator.New()

	// Register custom validations
	_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			match, _ := regexp.MatchString("^[a-z0-9]+(?:-[a-z0-9]+)*$", fl.Field().String())
			return match
		}
		return true
	})

	return &Validator{
		Validate: v,
	}
}

// GenerateValidationMessage generates a human-readable validation message
func (v *Validator) GenerateValidationMessage(field string, rule string) string {
	switch rule {
	case "required":
		return fmt.Sprintf("The %s field is required.", field)
	case "slug":
		return fmt.Sprintf("The %s must be a valid slug format.", field)
	case "min":
		return fmt.Sprintf("The %s field must be at least the minimum length.", field)
	case "max":
		return fmt.Sprintf("The %s field exceeds the maximum length.", field)
	default:
		return fmt.Sprintf("The %s field is invalid.", field)
	}
}

// GenerateValidationErrors converts validator errors to a slice of ValidationError
func (v *Validator) GenerateValidationErrors(err error) []ValidationError {
	var validations []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, value := range validationErrors {
			field, rule := value.Field(), value.Tag()
			validation := ValidationError{
				Field:   field,
				Message: v.GenerateValidationMessage(field, rule),
			}
			validations = append(validations, validation)
		}
	}

	return validations
}
