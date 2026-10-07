package repository

import (
	"github.com/OneDpsss/consultation-service/internal/models"

	"gorm.io/gorm"
)

// SlotRepository — слой доступа к данным для слотов консультаций
// преподавателя.
type SlotRepository struct {
	db *gorm.DB
}

// NewSlotRepository создаёт SlotRepository поверх переданного подключения
// GORM.
func NewSlotRepository(db *gorm.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

// CreateSlot создаёт новый слот консультации для преподавателя.
// Перед вставкой проверяется, что слот не пересекается по времени с уже
// существующими слотами этого же преподавателя.
func (r *SlotRepository) CreateSlot(slot *models.ConsultationSlot) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var overlapCount int64
		err := tx.Model(&models.ConsultationSlot{}).
			Where("teacher_id = ? AND starts_at < ? AND ends_at > ?",
				slot.TeacherID, slot.EndsAt, slot.StartsAt).
			Count(&overlapCount).Error
		if err != nil {
			return err
		}
		if overlapCount > 0 {
			return ErrSlotOverlap
		}
		return tx.Create(slot).Error
	})
}

// ListByTeacher возвращает все слоты консультаций, открытые
// преподавателем teacherID.
func (r *SlotRepository) ListByTeacher(teacherID uint) ([]models.ConsultationSlot, error) {
	var slots []models.ConsultationSlot
	err := r.db.Where("teacher_id = ?", teacherID).Find(&slots).Error
	return slots, err
}
