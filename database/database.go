package database

import (
	"context"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	once    sync.Once
	dialect = "sqlite"
	dsn     = "file:my_airline.db?cache=shared&_foreign_keys=1"
)

func InitDatabase(ctx context.Context) {
	once.Do(func() {
		switch dialect {
		case "sqlite":
			conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Printf("connect to SQLite err: %v \n", err)
				return
			}
			db = conn

		default:
			log.Println("unsupported dialect: ", dialect)
		}
	})

	log.Println("database connected successfully")
	return
}

func FinalizeDatabase(ctx context.Context) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to get sql.DB from gorm:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("fail to close database:", err)
	}

	log.Println("Database connection closed.")
}

// GetDB return global *gorm.DB
func GetDB() *gorm.DB {
	return db
}
