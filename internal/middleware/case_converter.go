package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

// bufferPool reuses bytes.Buffer instances to reduce allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// keyCache caches converted keys to avoid repeated conversions
type keyCache struct {
	snakeToCamel map[string]string
	camelToSnake map[string]string
	mu           sync.RWMutex
}

var cache = &keyCache{
	snakeToCamel: make(map[string]string),
	camelToSnake: make(map[string]string),
}

// getCachedSnakeToCamel returns cached camelCase conversion or computes and caches it
func getCachedSnakeToCamel(key string) string {
	cache.mu.RLock()
	if val, ok := cache.snakeToCamel[key]; ok {
		cache.mu.RUnlock()
		return val
	}
	cache.mu.RUnlock()

	// Compute conversion
	converted := strcase.ToLowerCamel(key)

	// Cache it (with size limit to prevent unbounded growth)
	cache.mu.Lock()
	if len(cache.snakeToCamel) < 10000 { // Limit cache size
		cache.snakeToCamel[key] = converted
	}
	cache.mu.Unlock()

	return converted
}

// getCachedCamelToSnake returns cached snake_case conversion or computes and caches it
func getCachedCamelToSnake(key string) string {
	cache.mu.RLock()
	if val, ok := cache.camelToSnake[key]; ok {
		cache.mu.RUnlock()
		return val
	}
	cache.mu.RUnlock()

	// Compute conversion
	converted := strcase.ToSnake(key)

	// Cache it (with size limit to prevent unbounded growth)
	cache.mu.Lock()
	if len(cache.camelToSnake) < 10000 { // Limit cache size
		cache.camelToSnake[key] = converted
	}
	cache.mu.Unlock()

	return converted
}

// CaseConverterMiddleware converts request body from camelCase to snake_case
// and response body from snake_case to camelCase
func CaseConverterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Early return for non-JSON content
		if !isJSONContent(c) {
			c.Next()
			return
		}

		// Process request body if present
		if err := processRequestBody(c); err != nil {
			// Log error but continue processing
			// You might want to handle this differently based on your requirements
			c.Next()
			return
		}

		// Create custom response writer
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bufferPool.Get().(*bytes.Buffer),
		}
		writer.body.Reset()
		c.Writer = writer

		// Process request
		c.Next()

		// Return buffer to pool when done
		defer func() {
			writer.body.Reset()
			bufferPool.Put(writer.body)
		}()

		// Only process JSON responses
		if !strings.Contains(writer.Header().Get("Content-Type"), "application/json") {
			// Write original response for non-JSON content
			if writer.status != 0 {
				writer.ResponseWriter.WriteHeader(writer.status)
			}
			writer.ResponseWriter.Write(writer.body.Bytes())
			return
		}

		// Process and write response
		processResponse(writer)
	}
}

// isJSONContent checks if the request should be processed
func isJSONContent(c *gin.Context) bool {
	// For GET requests, we only care about response processing
	if c.Request.Method == "GET" {
		return true
	}
	return strings.Contains(c.GetHeader("Content-Type"), "application/json")
}

// processRequestBody handles the request body conversion
func processRequestBody(c *gin.Context) error {
	// Skip if no body or GET request
	if c.Request.Body == nil || c.Request.Method == "GET" {
		return nil
	}

	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()

	if err != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return err
	}

	// Handle empty body
	if len(bodyBytes) == 0 {
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return nil
	}

	// Try to convert the body
	convertedBody, err := convertJSONKeys(bodyBytes, toSnakeCase)
	if err != nil {
		// Restore original body on error
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return nil // Don't propagate error, just use original
	}

	// Set converted body
	c.Request.Body = io.NopCloser(bytes.NewReader(convertedBody))
	c.Request.ContentLength = int64(len(convertedBody))
	return nil
}

// processResponse handles the response body conversion
func processResponse(writer *responseBodyWriter) {
	responseBody := writer.body.Bytes()

	// Handle empty response
	if len(responseBody) == 0 {
		if writer.status != 0 {
			writer.ResponseWriter.WriteHeader(writer.status)
		}
		return
	}

	// Try to convert response
	convertedBody, err := convertJSONKeys(responseBody, toCamelCase)

	// Write response (converted or original)
	if writer.status != 0 {
		writer.ResponseWriter.WriteHeader(writer.status)
	}

	if err != nil {
		writer.ResponseWriter.Write(responseBody)
	} else {
		writer.Header().Set("Content-Length", strconv.Itoa(len(convertedBody)))
		writer.ResponseWriter.Write(convertedBody)
	}
}

// convertJSONKeys converts JSON keys using the provided converter function
func convertJSONKeys(data []byte, converter func(string) string) ([]byte, error) {
	// Try to unmarshal as object first
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	// Convert keys recursively
	converted := convertValue(obj, converter)

	// Marshal back to JSON
	return json.Marshal(converted)
}

// convertValue recursively converts keys in a JSON value
func convertValue(v interface{}, converter func(string) string) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		return convertMap(val, converter)
	case []interface{}:
		return convertSlice(val, converter)
	default:
		return val
	}
}

// convertMap converts all keys in a map
func convertMap(m map[string]interface{}, converter func(string) string) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[converter(k)] = convertValue(v, converter)
	}
	return result
}

// convertSlice converts all values in a slice
func convertSlice(s []interface{}, converter func(string) string) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = convertValue(v, converter)
	}
	return result
}

// Case conversion functions (now using cache)
func toSnakeCase(s string) string {
	return getCachedCamelToSnake(s)
}

func toCamelCase(s string) string {
	return getCachedSnakeToCamel(s)
}

// responseBodyWriter captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	status  int
	written bool
}

// Write captures the response body
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	if !r.written {
		r.written = true
		if r.status == 0 {
			r.status = 200
		}
	}
	return r.body.Write(b)
}

// WriteHeader captures the status code
func (r *responseBodyWriter) WriteHeader(statusCode int) {
	if !r.written {
		r.status = statusCode
	}
}

// WriteString captures the response body
func (r *responseBodyWriter) WriteString(s string) (int, error) {
	if !r.written {
		r.written = true
		if r.status == 0 {
			r.status = 200
		}
	}
	return r.body.WriteString(s)
}
