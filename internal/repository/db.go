package repository

import (
	"consultation-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB opens a connection and runs auto-migration for all entities.
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Discipline{},
		&models.ConsultationSlot{},
		&models.Booking{},
		&models.Notification{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
