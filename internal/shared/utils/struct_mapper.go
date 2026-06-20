package utils

import (
	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// MapStructToUpdates dynamically maps struct fields to database updates
// If a field named "password" is found, it will be hashed before adding to updates
func MapStructToUpdates(req interface{}) (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	val := reflect.ValueOf(req)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		jsonTag := typ.Field(i).Tag.Get("json")

		// Skip empty fields and get the actual field name from JSON tag
		if field.String() != "" {
			// Remove omitempty suffix if present
			if strings.Contains(jsonTag, ",") {
				jsonTag = strings.Split(jsonTag, ",")[0]
			}

			// Check if this is a password field and hash it
			if strings.ToLower(fieldName) == "password" {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(field.String()), bcrypt.DefaultCost)
				if err != nil {
					return nil, err
				}
				updates[jsonTag] = string(hashedPassword)
			} else {
				updates[jsonTag] = field.String()
			}
		}
	}

	return updates, nil
}
