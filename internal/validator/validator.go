package validators

import (
	"fmt"
	"regexp"
	"strings"

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

// GenerateValidationMessage generates a human-readable validation message (Laravel-style)
func (v *Validator) GenerateValidationMessage(field string, rule string, param string) string {
	// Convert field to snake_case for display
	fieldName := strings.ToLower(strings.ReplaceAll(field, " ", "_"))

	switch rule {
	case "required":
		return fmt.Sprintf("The %s field is required.", fieldName)
	case "email":
		return fmt.Sprintf("The %s must be a valid email address.", fieldName)
	case "slug":
		return fmt.Sprintf("The %s must be a valid slug format.", fieldName)
	case "min":
		return fmt.Sprintf("The %s must be at least %s characters.", fieldName, param)
	case "max":
		return fmt.Sprintf("The %s may not be greater than %s characters.", fieldName, param)
	default:
		return fmt.Sprintf("The %s field is invalid.", fieldName)
	}
}

// GenerateValidationErrors converts validator errors to a slice of ValidationError
func (v *Validator) GenerateValidationErrors(err error) []ValidationError {
	var validations []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, value := range validationErrors {
			field := value.Field()
			rule := value.Tag()
			param := ""

			// Get parameter for rules like min, max
			if value.Param() != "" {
				param = value.Param()
			}

			validation := ValidationError{
				Field:   field,
				Message: v.GenerateValidationMessage(field, rule, param),
			}
			validations = append(validations, validation)
		}
	}

	return validations
}
