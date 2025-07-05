package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a GORM database connection object
// with a configured connection pool.
func InitDB() (*gorm.DB, error) {
	// --- 1. Database Connection String (DSN) ---
	// It's good practice to fetch these values from environment variables
	// or a configuration file in a real application.
	dsn := "host=localhost user=root password=11111111 dbname=lotto_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@lotto_db"

	// --- 2. Open GORM Database Connection ---
	// gorm.Open creates an initial connection and prepares a connection pool.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// --- 3. Get the Underlying sql.DB object ---
	// GORM uses the standard library's database/sql package for connection pooling.
	// We need to get this underlying object to configure the pool.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// --- 4. Configure the Connection Pool ---
	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool. If n <= 0, no idle connections are retained.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's age.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// --- 5. Verify the Connection ---
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Successfully connected to the database with a configured connection pool!")
	return db, nil
}
