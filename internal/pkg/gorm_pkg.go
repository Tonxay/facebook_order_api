package gormpkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// gormdb is a global singleton variable for gorm.DB
var gormdb *gorm.DB

// Init initializes gormdb
func Init(webhook string) error {
	var err error
	if webhook != "webhook" {
		// Load .env file if exists
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, reading from system environment")
		}
	}

	DBURL := os.Getenv("PSQL_DB_URL")
	pool, err := pgxpool.New(context.Background(), DBURL)

	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %w", err)
	}

	log.Println("✅ Connected to the database successfully")

	gormdb, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return nil
}

// GetDB returns gormdb
func GetDB() *gorm.DB {
	if gormdb == nil {
		Init("")
	}
	return gormdb
}
