package database

import (
	"log"
	"time"

	"replme/model"
	"replme/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(pgUrl string) {
	var err error
	for i := 0; i < 10; i++ {
		util.SLogger.Infof("Connecting to DB, try %d", i)
		DB, err = gorm.Open(postgres.Open(pgUrl), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
}

func Migrate() {
	DB.AutoMigrate(&model.User{}, &model.Devenv{})
}
