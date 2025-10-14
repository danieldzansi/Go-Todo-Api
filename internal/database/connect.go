// package database

// import (
// 	"fmt"
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func ConnectGorm() (*gorm.DB, error) {
// 	host := os.Getenv("DB_HOST")
// 	if host == "" {
// 		host = ("localhost")
// 	}
// 	port := os.Getenv("DB_PORT")
// 	if port == "" {
// 		port = ("5432")
// 	}
// 	user := os.Getenv("DB_USER")
// 	if user == "" {
// 		user = ("postgres")
// 	}
// 	password := os.Getenv("DB_PASSWORD")
// 	if password == "" {
// 		password = "2345"
// 	}
// 	dbname := os.Getenv("DB_NAME")
// 	if dbname == "" {
// 		dbname = ("trial_db")
// 	}

// 	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to open database connection")
// 	}
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
// 	}
// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, fmt.Errorf("failed to ping database: %w", err)
// 	}
// 	sqlDB.SetMaxOpenConns(25)
// 	sqlDB.SetMaxIdleConns(25)

// 	return db, nil
// }

// package database

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func ConnectGorm() (*gorm.DB, error) {
// 	// Try to load connection from Render DATABASE_URL
// 	dsn := os.Getenv("DATABASE_URL")
// 	if dsn == "" {
// 		// Local fallback (for running on your laptop)
// 		host := os.Getenv("DB_HOST")
// 		if host == "" {
// 			host = "localhost"
// 		}
// 		port := os.Getenv("DB_PORT")
// 		if port == "" {
// 			port = "5432"
// 		}
// 		user := os.Getenv("DB_USER")
// 		if user == "" {
// 			user = "postgres"
// 		}
// 		password := os.Getenv("DB_PASSWORD")
// 		if password == "" {
// 			password = "2345"
// 		}
// 		dbname := os.Getenv("DB_NAME")
// 		if dbname == "" {
// 			dbname = "trial_db"
// 		}

// 		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 			host, port, user, password, dbname)
// 	}

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
// 	}

// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, fmt.Errorf("failed to ping database: %w", err)
// 	}

// 	sqlDB.SetMaxOpenConns(25)
// 	sqlDB.SetMaxIdleConns(25)

// 	log.Println("✅ Connected to PostgreSQL successfully!")

// 	return db, nil
// }

package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGorm() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn != "" && !strings.Contains(dsn, "sslmode") {
		dsn += "?sslmode=require"
	} else if dsn == "" {
		// Local fallback
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "2345"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "trial_db"
		}

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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

	log.Println("✅ Connected to Render PostgreSQL successfully!")
	return db, nil
}
