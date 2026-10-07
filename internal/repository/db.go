// Package repository содержит слой доступа к данным: подключение к
// PostgreSQL через GORM и репозитории для слотов и записей на консультацию.
package repository

import (
	"consultation-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB opens a PostgreSQL connection using dsn and runs auto-migration
// for every entity in internal/models, creating or updating tables as
// needed. Call this once at startup (see cmd/server/main.go).
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
