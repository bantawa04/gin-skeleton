package transformer

import (
	"reflect"
)

// TransformFunc is a function that transforms one type to another
type TransformFunc[T, U any] func(T) U

// TransformCollection transforms a slice of one type to a slice of another type using the provided transform function
func TransformCollection[T, U any](items []T, transformFn TransformFunc[T, U]) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = transformFn(item)
	}
	return result
}

// TransformPagination transforms a paginated result by applying the transform function to the data field
func TransformPagination[T, U any](paginatedResult map[string]interface{}, transformFn TransformFunc[T, U]) map[string]interface{} {
	// Create a copy of the paginated result to avoid modifying the original
	result := make(map[string]interface{})
	for k, v := range paginatedResult {
		if k == "data" {
			// Handle data field separately
			continue
		}
		result[k] = v
	}

	// Check if the data field exists and is a slice of the expected type
	if data, ok := paginatedResult["data"]; ok {
		// Use reflection to check if data is a slice of T
		dataValue := reflect.ValueOf(data)
		if dataValue.Kind() == reflect.Slice {
			// Try to convert to []T
			if items, ok := data.([]T); ok {
				// Transform the items and update the data field
				result["data"] = TransformCollection(items, transformFn)
			} else {
				// If we can't convert, just copy the original data
				result["data"] = data
			}
		} else {
			// If it's not a slice, just copy it
			result["data"] = data
		}
	}

	return result
}
