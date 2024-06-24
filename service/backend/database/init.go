package database

import (
	"log"
	"os"

	"replme/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	sqlitePath := os.Getenv("REPL_SQLITE")

	if sqlitePath == "" {
		log.Fatal("No Sqlitepath")
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
}

func Migrate() {
	DB.AutoMigrate(&model.User{}, &model.Devenv{})
}
