package utils

import (
	validators "gin/internal/validator"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

// ExtractBindingErrors converts Gin binding errors to ValidationError format
func ExtractBindingErrors(err error) []validators.ValidationError {
	var validationErrors []validators.ValidationError

	// Check if it's a validator.ValidationErrors (from struct validation)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrs {
			field := ve.Field()
			rule := ve.Tag()
			param := ve.Param()

			// Generate readable message
			message := generateBindingErrorMessage(field, rule, param)

			validationErrors = append(validationErrors, validators.ValidationError{
				Field:   field,
				Message: message,
			})
		}
		return validationErrors
	}

	// For other binding errors, try to extract field names from error message
	errorStr := err.Error()

	// Check for type mismatch errors (e.g., "email": 1 when string expected)
	if strings.Contains(errorStr, "cannot unmarshal") {
		// Try to extract field name from type mismatch error
		fieldName := extractFieldFromTypeMismatchError(errorStr)
		if fieldName != "" {
			validationErrors = append(validationErrors, validators.ValidationError{
				Field:   fieldName,
				Message: generateTypeMismatchMessage(fieldName, errorStr),
			})
			return validationErrors
		}
	}

	// Check for JSON syntax errors
	if strings.Contains(errorStr, "invalid character") || strings.Contains(errorStr, "invalid syntax") {
		// This is a JSON parsing error, not a field-specific error
		// Return empty to let the handler show a generic JSON syntax error
		return validationErrors
	}

	// Parse multiple field errors from the error string
	// Pattern: "Key: 'Request.Field1' Error:...\nKey: 'Request.Field2' Error:..."
	lines := strings.Split(errorStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Key: '") {
			fieldName := extractFieldFromError(line)
			if fieldName != "" {
				// Determine the rule and param from the error message
				rule, param := extractRuleAndParam(line)

				message := generateBindingErrorMessage(fieldName, rule, param)
				if rule == "type_mismatch" {
					message = generateTypeMismatchMessage(fieldName, line)
				}

				validationErrors = append(validationErrors, validators.ValidationError{
					Field:   fieldName,
					Message: message,
				})
			}
		}
	}

	return validationErrors
}

// generateBindingErrorMessage creates a readable error message
func generateBindingErrorMessage(field, rule, param string) string {
	// Convert field name to lower camel case for display (e.g., FirstName -> firstName)
	fieldName := strcase.ToLowerCamel(field)

	switch rule {
	case "required":
		return "The " + fieldName + " field is required."
	case "email":
		return "The " + fieldName + " must be a valid email address."
	case "min":
		if param != "" {
			return "The " + fieldName + " must be at least " + param + " characters."
		}
		return "The " + fieldName + " field is too short."
	case "max":
		if param != "" {
			return "The " + fieldName + " may not be greater than " + param + " characters."
		}
		return "The " + fieldName + " field is too long."
	default:
		return "The " + fieldName + " field is invalid."
	}
}


// extractFieldFromError tries to extract field name from error message
func extractFieldFromError(errorMsg string) string {
	var fieldName string

	// Look for pattern: Field validation for 'FieldName'
	if strings.Contains(errorMsg, "Field validation for '") {
		start := strings.Index(errorMsg, "Field validation for '") + len("Field validation for '")
		end := strings.Index(errorMsg[start:], "'")
		if end > 0 {
			fieldName = errorMsg[start : start+end]
		}
	}

	// If not found, look for pattern: Key: 'Request.FieldName'
	if fieldName == "" && strings.Contains(errorMsg, "Key: '") {
		parts := strings.Split(errorMsg, "Key: '")
		if len(parts) > 1 {
			keyPart := parts[1]
			if dotIndex := strings.Index(keyPart, "."); dotIndex > 0 {
				fieldPart := keyPart[dotIndex+1:]
				if endIndex := strings.Index(fieldPart, "'"); endIndex > 0 {
					fieldName = fieldPart[:endIndex]
				}
			}
		}
	}

	return fieldName
}

// extractRuleAndParam extracts validation rule and parameter from error message
func extractRuleAndParam(errorMsg string) (rule, param string) {
	// Check for min rule with parameter (e.g., "min=6")
	if strings.Contains(errorMsg, "min") {
		rule = "min"
		// Try to extract the parameter value
		if minIndex := strings.Index(errorMsg, "min="); minIndex > 0 {
			paramStart := minIndex + len("min=")
			paramEnd := paramStart
			for paramEnd < len(errorMsg) && (errorMsg[paramEnd] >= '0' && errorMsg[paramEnd] <= '9') {
				paramEnd++
			}
			if paramEnd > paramStart {
				param = errorMsg[paramStart:paramEnd]
			}
		}
		return
	}

	// Check for max rule with parameter (e.g., "max=100")
	if strings.Contains(errorMsg, "max") {
		rule = "max"
		// Try to extract the parameter value
		if maxIndex := strings.Index(errorMsg, "max="); maxIndex > 0 {
			paramStart := maxIndex + len("max=")
			paramEnd := paramStart
			for paramEnd < len(errorMsg) && (errorMsg[paramEnd] >= '0' && errorMsg[paramEnd] <= '9') {
				paramEnd++
			}
			if paramEnd > paramStart {
				param = errorMsg[paramStart:paramEnd]
			}
		}
		return
	}

	// Check for other rules
	if strings.Contains(errorMsg, "required") || strings.Contains(errorMsg, "Field validation for") {
		rule = "required"
	} else if strings.Contains(errorMsg, "email") {
		rule = "email"
	} else if strings.Contains(errorMsg, "cannot unmarshal") {
		rule = "type_mismatch"
	} else {
		rule = "invalid"
	}

	return
}

// extractFieldFromTypeMismatchError extracts field name from type mismatch errors
func extractFieldFromTypeMismatchError(errorMsg string) string {
	// Pattern: "json: cannot unmarshal number into Go struct field SignupRequest.email of type string"
	// Pattern: "json: cannot unmarshal bool into Go struct field SignupRequest.email of type string"
	// Pattern: "json: cannot unmarshal number into Go value of type string"

	// Look for pattern: "field Request.FieldName" or "Request.FieldName"
	if strings.Contains(errorMsg, "field ") {
		parts := strings.Split(errorMsg, "field ")
		if len(parts) > 1 {
			fieldPart := parts[1]
			if dotIndex := strings.Index(fieldPart, "."); dotIndex > 0 {
				fieldName := fieldPart[dotIndex+1:]
				// Remove " of type" or anything after space
				if spaceIndex := strings.Index(fieldName, " "); spaceIndex > 0 {
					fieldName = fieldName[:spaceIndex]
				}
				// Remove trailing period if present
				fieldName = strings.TrimSuffix(fieldName, ".")
				return fieldName
			}
		}
	}

	// Alternative pattern: look for struct field directly
	// "json: cannot unmarshal number into Go struct field SignupRequest.email"
	if strings.Contains(errorMsg, "struct field ") {
		parts := strings.Split(errorMsg, "struct field ")
		if len(parts) > 1 {
			fieldPart := parts[1]
			if dotIndex := strings.Index(fieldPart, "."); dotIndex > 0 {
				fieldName := fieldPart[dotIndex+1:]
				// Remove " of type" or anything after space
				if spaceIndex := strings.Index(fieldName, " "); spaceIndex > 0 {
					fieldName = fieldName[:spaceIndex]
				}
				// Remove trailing period if present
				fieldName = strings.TrimSuffix(fieldName, ".")
				return fieldName
			}
		}
	}

	return ""
}

// generateTypeMismatchMessage creates a readable message for type mismatch errors
func generateTypeMismatchMessage(fieldName, errorMsg string) string {
	fieldNameSnake := strings.ToLower(strcase.ToSnake(fieldName))

	// Try to determine expected type from error message
	expectedType := "string"
	if strings.Contains(errorMsg, "of type string") {
		expectedType = "string"
	} else if strings.Contains(errorMsg, "of type int") || strings.Contains(errorMsg, "of type number") {
		expectedType = "number"
	} else if strings.Contains(errorMsg, "of type bool") {
		expectedType = "boolean"
	}

	// Try to determine actual type from error message
	actualType := "invalid type"
	if strings.Contains(errorMsg, "cannot unmarshal number") {
		actualType = "number"
	} else if strings.Contains(errorMsg, "cannot unmarshal bool") {
		actualType = "boolean"
	} else if strings.Contains(errorMsg, "cannot unmarshal string") {
		actualType = "string"
	}

	if actualType != "invalid type" && expectedType != actualType {
		return "The " + fieldNameSnake + " field must be a " + expectedType + ", but received " + actualType + "."
	}

	return "The " + fieldNameSnake + " field must be a " + expectedType + "."
}
