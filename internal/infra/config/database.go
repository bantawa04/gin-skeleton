package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetries        = 5
	initialDelay      = 1 * time.Second
	maxDelay          = 30 * time.Second
	backoffMultiplier = 2
)

// InitDatabase establishes a connection to the database using the application config
// with retry logic and exponential backoff
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dbConfig := cfg.Database()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User,
		dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Retry connection with exponential backoff
	var db *gorm.DB
	var err error
	delay := initialDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			// Connection successful
			break
		}

		if attempt < maxRetries {
			log.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...",
				attempt, maxRetries, err, delay)
			time.Sleep(delay)

			// Exponential backoff with max delay cap
			delay *= backoffMultiplier
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Verify connection is alive
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return db, nil
}
