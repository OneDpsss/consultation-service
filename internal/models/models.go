// @file models.go
// @brief Модель предметной области системы записи на консультации.
// @copyright 2026
// @author Шавлягин Валерий, Морехон Энрике
//
// Содержит классы предметной области: User (и его роли), Department,
// Discipline, ConsultationSlot, Booking, Notification — ровно та модель,
// что описана в Приложении Б технического задания.

package models

import "time"

// Role — роль пользователя в системе.
type Role string

const (
	RoleStudent Role = "student" ///< студент
	RoleTeacher Role = "teacher" ///< преподаватель
	RoleAdmin   Role = "admin"   ///< администратор
)

// / @brief Обобщённая сущность пользователя системы.
// /
// / Содержит данные для входа в систему. Student и Teacher не выделены
// / отдельными Go-типами (в отличие от UML-модели в ТЗ), а представлены
// / полем Role — это осознанное упрощение для реляционного хранения:
// / единая таблица users вместо наследования таблиц.
type User struct {
	ID           uint      `gorm:"primaryKey"`           ///< идентификатор пользователя
	Email        string    `gorm:"uniqueIndex;not null"` ///< логин (email)
	PasswordHash string    `gorm:"not null"`             ///< хеш пароля
	FullName     string    `gorm:"not null"`             ///< ФИО
	Role         Role      `gorm:"not null"`             ///< роль: студент/преподаватель/администратор
	CreatedAt    time.Time ///< дата регистрации
}

// / @brief Кафедра образовательной организации.
type Department struct {
	ID   uint   `gorm:"primaryKey"` ///< идентификатор кафедры
	Name string `gorm:"not null"`   ///< наименование кафедры
}

// / @brief Учебная дисциплина, привязанная к кафедре.
type Discipline struct {
	ID           uint       `gorm:"primaryKey"` ///< идентификатор дисциплины
	Name         string     `gorm:"not null"`   ///< наименование дисциплины
	DepartmentID uint       ///< внешний ключ на кафедру
	Department   Department `gorm:"foreignKey:DepartmentID"` ///< кафедра дисциплины
}

// / @brief Консультационный слот — временной интервал, выделенный
// / преподавателем для консультации по дисциплине.
// /
// / Связан с Преподавателем, Дисциплиной, датой и временем проведения —
// / как и описано в модели предметной области ТЗ (Приложение Б).
type ConsultationSlot struct {
	ID           uint       `gorm:"primaryKey"`           ///< идентификатор слота
	TeacherID    uint       `gorm:"not null;index"`       ///< внешний ключ на преподавателя
	Teacher      User       `gorm:"foreignKey:TeacherID"` ///< преподаватель, создавший слот
	DisciplineID uint       ///< внешний ключ на дисциплину
	Discipline   Discipline `gorm:"foreignKey:DisciplineID"` ///< дисциплина консультации
	StartsAt     time.Time  `gorm:"not null"`                ///< время начала консультации
	EndsAt       time.Time  `gorm:"not null"`                ///< время окончания консультации
	Capacity     int        `gorm:"not null;default:1"`      ///< вместимость слота (число студентов)
}

// BookingStatus — статус записи на консультацию.
type BookingStatus string

const (
	BookingActive    BookingStatus = "active"    ///< запись активна
	BookingCancelled BookingStatus = "cancelled" ///< запись отменена студентом
	BookingDone      BookingStatus = "done"      ///< консультация проведена
)

// / @brief Запись студента на консультацию.
// /
// / Связывает Студента с конкретной Консультацией и содержит статус
// / записи — формирование уведомления происходит при изменении этого
// / статуса (см. Notification).
type Booking struct {
	ID        uint             `gorm:"primaryKey"`              ///< идентификатор записи
	StudentID uint             `gorm:"not null;index"`          ///< внешний ключ на студента
	Student   User             `gorm:"foreignKey:StudentID"`    ///< студент, сделавший запись
	SlotID    uint             `gorm:"not null;index"`          ///< внешний ключ на слот
	Slot      ConsultationSlot `gorm:"foreignKey:SlotID"`       ///< слот консультации
	Status    BookingStatus    `gorm:"not null;default:active"` ///< статус записи
	CreatedAt time.Time        ///< дата создания записи
}

// / @brief Уведомление пользователя о событии, связанном с записью или
// / консультационным слотом.
// /
// / Формируется при изменении состояния записи или слота (создание,
// / изменение, отмена) — требование п. 3.1.1 ТЗ.
type Notification struct {
	ID        uint      `gorm:"primaryKey"`             ///< идентификатор уведомления
	UserID    uint      `gorm:"not null;index"`         ///< внешний ключ на адресата
	Message   string    `gorm:"not null"`               ///< текст уведомления
	IsRead    bool      `gorm:"not null;default:false"` ///< признак прочтения
	CreatedAt time.Time ///< дата создания уведомления
}
