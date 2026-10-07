// Package repository — слой доступа к данным: подключение к PostgreSQL
// через GORM и репозитории для слотов и записей на консультации.
package repository

import (
	"consultation-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB открывает соединение с PostgreSQL по dsn и прогоняет
// автомиграцию для всех сущностей из internal/models, создавая или
// обновляя таблицы. Вызывается один раз при старте (см. cmd/server/main.go).
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
