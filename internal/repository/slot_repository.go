package repository

import (
	"consultation-service/internal/models"

	"gorm.io/gorm"
)

type SlotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

// CreateSlot inserts a new consultation slot for a teacher.
// The overlap check (ТЗ 3.2.2 — "отсутствие пересечений консультационных
// слотов у преподавателя") is enforced before insertion.
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

func (r *SlotRepository) ListByTeacher(teacherID uint) ([]models.ConsultationSlot, error) {
	var slots []models.ConsultationSlot
	err := r.db.Where("teacher_id = ?", teacherID).Find(&slots).Error
	return slots, err
}
