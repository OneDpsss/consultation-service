// Package models описывает доменную модель сервиса записи на консультации:
// пользователей и их роли, кафедры, дисциплины, слоты консультаций, записи
// студентов и уведомления.
package models

import "time"

// Role — роль пользователя в системе.
type Role string

const (
	RoleStudent Role = "student" // студент
	RoleTeacher Role = "teacher" // преподаватель
	RoleAdmin   Role = "admin"   // администратор
)

// User — базовая сущность с данными для входа в систему.
//
// Студент и преподаватель не выделены в отдельные типы — роль хранится
// в поле Role. Так проще: одна таблица users вместо отдельных таблиц
// на каждую роль.
type User struct {
	ID           uint
	Email        string `gorm:"uniqueIndex;not null"` // логин (email)
	PasswordHash string `gorm:"not null"`             // хеш пароля
	FullName     string `gorm:"not null"`             // ФИО
	Role         Role   `gorm:"not null"`             // роль пользователя
	CreatedAt    time.Time
}

// TeacherProfile расширяет User данными, специфичными для преподавателя
// (кафедра). Встраивание User даёт TeacherProfile все его поля плюс свои.
type TeacherProfile struct {
	User
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

// Department — кафедра.
type Department struct {
	ID   uint
	Name string `gorm:"not null"`
}

// Discipline — учебная дисциплина, привязанная к кафедре.
type Discipline struct {
	ID           uint
	Name         string `gorm:"not null"`
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

// ConsultationSlot — временной интервал, который преподаватель открывает
// для консультации по дисциплине.
type ConsultationSlot struct {
	ID           uint
	TeacherID    uint `gorm:"not null;index"`
	Teacher      User `gorm:"foreignKey:TeacherID"`
	DisciplineID uint
	Discipline   Discipline `gorm:"foreignKey:DisciplineID"`
	StartsAt     time.Time  `gorm:"not null"`
	EndsAt       time.Time  `gorm:"not null"`
	Capacity     int        `gorm:"not null;default:1"` // сколько студентов вмещает слот
}

// BookingStatus — статус записи студента на консультацию.
type BookingStatus string

const (
	BookingActive    BookingStatus = "active"    // запись активна
	BookingCancelled BookingStatus = "cancelled" // запись отменена студентом
	BookingDone      BookingStatus = "done"      // консультация проведена
)

// Booking связывает студента с конкретной консультацией и хранит статус
// записи. При смене статуса создаётся Notification.
type Booking struct {
	ID        uint
	StudentID uint             `gorm:"not null;index"`
	Student   User             `gorm:"foreignKey:StudentID"`
	SlotID    uint             `gorm:"not null;index"`
	Slot      ConsultationSlot `gorm:"foreignKey:SlotID"`
	Status    BookingStatus    `gorm:"not null;default:active"`
	CreatedAt time.Time
}

// Notification — уведомление пользователю о событии, связанном с записью
// или слотом (создание, изменение, отмена).
type Notification struct {
	ID        uint
	UserID    uint   `gorm:"not null;index"`
	Message   string `gorm:"not null"`
	IsRead    bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
}
