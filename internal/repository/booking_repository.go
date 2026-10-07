package repository

import (
	"errors"

	"consultation-service/internal/models"

	"gorm.io/gorm"
)

var (
	// ErrSlotFull is returned by CreateBooking when the slot's capacity
	// is already reached (ТЗ 3.2.1 — запрет записи двух студентов на один
	// и тот же индивидуальный слот).
	ErrSlotFull = errors.New("slot has no free capacity")

	// ErrSlotOverlap is returned by SlotRepository.CreateSlot when the new
	// slot's time range overlaps an existing slot of the same teacher.
	ErrSlotOverlap = errors.New("slot overlaps with an existing one for this teacher")
)

// BookingRepository is the data-access layer for student bookings.
type BookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a BookingRepository over the given GORM
// connection.
func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBooking books slotID for studentID.
//
// The capacity check and insert run inside a single transaction with a
// row-level lock (FOR UPDATE) on the slot, so two concurrent requests
// cannot both pass the check and overbook the same slot — this is what
// enforces ТЗ 3.2.1 at the database level, not just in application code.
// Returns ErrSlotFull if the slot has no free capacity.
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

// CancelBooking marks bookingID as cancelled, but only if it belongs to
// studentID — the WHERE clause doubles as an ownership check, so one
// student cannot cancel another student's booking by guessing an id.
//
// It is not an error if no row matches (wrong id, wrong owner, or an
// already-cancelled booking): GORM reports 0 rows affected without
// returning an error, so the caller gets a nil error either way. The
// handler layer treats this as success, matching ТЗ 3.1.1 "Отменить
// запись на консультацию".
func (r *BookingRepository) CancelBooking(bookingID, studentID uint) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ? AND student_id = ?", bookingID, studentID).
		Update("status", models.BookingCancelled).Error
}

// ListByStudent returns every booking made by studentID, each with its
// ConsultationSlot preloaded, so callers get the slot's time/discipline
// data without a second query. Used by the "просмотр студентом списка
// своих записей" use case (ТЗ 3.1.1).
func (r *BookingRepository) ListByStudent(studentID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Slot").Where("student_id = ?", studentID).Find(&bookings).Error
	return bookings, err
}
