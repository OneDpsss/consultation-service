package repository

import (
	"errors"

	"github.com/OneDpsss/consultation-service/internal/models"

	"gorm.io/gorm"
)

var (
	// ErrSlotFull возвращается из CreateBooking, когда в слоте больше нет
	// свободных мест.
	ErrSlotFull = errors.New("slot has no free capacity")

	// ErrSlotOverlap возвращается из SlotRepository.CreateSlot, если новый
	// слот по времени пересекается с уже существующим слотом того же
	// преподавателя.
	ErrSlotOverlap = errors.New("slot overlaps with an existing one for this teacher")
)

// BookingRepository — слой доступа к данным для записей студентов.
type BookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository создаёт BookingRepository поверх переданного
// подключения GORM.
func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBooking записывает студента studentID на слот slotID.
//
// Проверка вместимости и вставка выполняются в одной транзакции с
// блокировкой строки слота (FOR UPDATE), поэтому два одновременных
// запроса не смогут оба пройти проверку и переполнить один и тот же
// слот. Если свободных мест нет, возвращает ErrSlotFull.
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

// CancelBooking отмечает запись bookingID как отменённую, но только если
// она принадлежит studentID — условие в WHERE заодно проверяет
// владельца, так что один студент не может отменить чужую запись,
// подобрав id.
//
// Если ни одна строка не подошла (неверный id, чужая запись или запись
// уже отменена) — это не считается ошибкой: GORM в таком случае просто
// сообщает 0 затронутых строк, err остаётся nil.
func (r *BookingRepository) CancelBooking(bookingID, studentID uint) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ? AND student_id = ?", bookingID, studentID).
		Update("status", models.BookingCancelled).Error
}

// ListByStudent возвращает все записи студента studentID вместе с
// предзагруженным слотом (Slot), чтобы не делать по нему отдельный
// запрос.
func (r *BookingRepository) ListByStudent(studentID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Slot").Where("student_id = ?", studentID).Find(&bookings).Error
	return bookings, err
}
