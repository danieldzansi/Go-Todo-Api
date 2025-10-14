package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGorm() (*gorm.DB, error) {
	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = ("localhost")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = ("5432")
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = ("postgres")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "2345"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = ("trial_db")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to open database connection")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	return db, nil
}
