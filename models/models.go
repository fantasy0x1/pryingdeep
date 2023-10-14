package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/r00tk3y/prying-deep/pkg/logger"
)

var db *gorm.DB

type Model struct {
	ID uint `gorm:"primaryKey"`
}

func SetupDatabase(dbUrl string) {
	var err error
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		logger.Fatalf("models.Setup err: %v", err)
	}

	err = db.AutoMigrate(&WebPage{}, &WordpressFootPrint{}, &Email{}, &PhoneNumber{})
	if err != nil {
		logger.Errorf("error during AutoMigrations", err)
	}

}

// DeleteAllRecords This function is primarily needed for clearing test database,
// but there might be support later on for clearing out all records
