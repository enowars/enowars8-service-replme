package database

import (
	"log"

	"replme/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	DB.Raw("PRAGMA synchronous = NORMAL;")
	DB.Raw("PRAGMA journal_mode = WAL;")
	DB.Raw("PRAGMA temp_store = MEMORY;")
	DB.Raw("PRAGMA cache_size = 10000;")
	DB.Raw("PRAGMA mmap_size = 268435456;")
	DB.Raw("PRAGMA optimize;")
}

func Migrate() {
	DB.AutoMigrate(&model.User{}, &model.Devenv{})
}
