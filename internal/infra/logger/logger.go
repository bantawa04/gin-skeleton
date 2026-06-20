package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	// Logger is the global logger instance
	Logger *logrus.Logger

	// ErrorLogger is specifically for errors (status > 400)
	ErrorLogger *logrus.Logger
)

// Init initializes the logger with file rotation and cleanup
func Init() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Initialize main logger
	Logger = logrus.New()

	// Initialize error logger
	ErrorLogger = logrus.New()

	// Configure main logger
	if err := configureLogger(Logger, filepath.Join(logsDir, "app.log")); err != nil {
		return fmt.Errorf("failed to configure main logger: %w", err)
	}

	// Configure error logger
	if err := configureLogger(ErrorLogger, filepath.Join(logsDir, "errors.log")); err != nil {
		return fmt.Errorf("failed to configure error logger: %w", err)
	}

	// Set error logger to only log errors and above
	ErrorLogger.SetLevel(logrus.ErrorLevel)

	return nil
}

// configureLogger sets up a logger with file rotation and cleanup
func configureLogger(logger *logrus.Logger, logPath string) error {
	// Create file rotator with date-based naming
	// Format: app-2025-01-16.log, app-2025-01-17.log, etc.
	rotator, err := rotatelogs.New(
		logPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(logPath),           // Create symlink to current log
		rotatelogs.WithRotationTime(24*time.Hour),  // Rotate daily at midnight
		rotatelogs.WithMaxAge(30*24*time.Hour),     // Keep logs for 30 days
		rotatelogs.WithRotationSize(100*1024*1024), // Also rotate if file exceeds 100MB
	)
	if err != nil {
		return fmt.Errorf("failed to create rotator: %w", err)
	}

	// Set output to rotator
	logger.SetOutput(rotator)

	// Set formatter to JSON for better parsing
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	return nil
}

// LogHTTPError logs HTTP errors with status code > 400
func LogHTTPError(statusCode int, message string, fields map[string]interface{}) {
	if statusCode >= 400 {
		// Add status code to fields
		if fields == nil {
			fields = make(map[string]interface{})
		}
		fields["status_code"] = statusCode
		fields["type"] = "http_error"

		// Determine log level based on status code
		var level logrus.Level
		switch {
		case statusCode >= 500:
			level = logrus.ErrorLevel
		case statusCode >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		// Log to both loggers
		Logger.WithFields(fields).Log(level, message)
		ErrorLogger.WithFields(fields).Log(level, message)
	}
}

// LogError logs general errors
func LogError(err error, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	fields["type"] = "general_error"

	Logger.WithFields(fields).Error(message)
	ErrorLogger.WithFields(fields).Error(message)
}

// LogInfo logs informational messages
func LogInfo(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Info(message)
}

// LogWarning logs warning messages
func LogWarning(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Warn(message)
}

// LogDebug logs debug messages
func LogDebug(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Debug(message)
}

// GetLogger returns the main logger instance
func GetLogger() *logrus.Logger {
	return Logger
}

// GetErrorLogger returns the error logger instance
func GetErrorLogger() *logrus.Logger {
	return ErrorLogger
}
