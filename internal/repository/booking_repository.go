package repository

import (
	"errors"

	"consultation-service/internal/models"

	"gorm.io/gorm"
)

var (
	ErrSlotFull    = errors.New("slot has no free capacity")
	ErrSlotOverlap = errors.New("slot overlaps with an existing one for this teacher")
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBooking books a slot for a student, enforcing capacity inside a
// transaction so two concurrent requests cannot both pass the check (ТЗ 3.2.1).
func (r *BookingRepository) CreateBooking(studentID, slotID uint) (*models.Booking, error) {
	var booking models.Booking

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var slot models.ConsultationSlot
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&slot, slotID).Error; err != nil {
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Booking{}).
			Where("slot_id = ? AND status = ?", slotID, models.BookingActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= slot.Capacity {
			return ErrSlotFull
		}

		booking = models.Booking{
			StudentID: studentID,
			SlotID:    slotID,
			Status:    models.BookingActive,
		}
		return tx.Create(&booking).Error
	})

	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) CancelBooking(bookingID, studentID uint) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ? AND student_id = ?", bookingID, studentID).
		Update("status", models.BookingCancelled).Error
}

func (r *BookingRepository) ListByStudent(studentID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Slot").Where("student_id = ?", studentID).Find(&bookings).Error
	return bookings, err
}
