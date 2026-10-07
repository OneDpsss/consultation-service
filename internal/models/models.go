// Package models defines the domain model of the consultation booking
// system: User and its roles, Department, Discipline, ConsultationSlot,
// Booking, Notification — the same model described in Приложение Б
// технического задания.
package models

import "time"

// Role is a user role in the system.
type Role string

const (
	RoleStudent Role = "student" // студент
	RoleTeacher Role = "teacher" // преподаватель
	RoleAdmin   Role = "admin"   // администратор
)

// User is the generalized entity holding login data.
//
// Student and Teacher are not separate Go types (unlike the UML model in
// the ТЗ) — they are represented by the Role field instead. This is a
// deliberate simplification for relational storage: one users table
// instead of table-per-subtype inheritance.
type User struct {
	ID           uint   // идентификатор пользователя
	Email        string `gorm:"uniqueIndex;not null"` // логин (email)
	PasswordHash string `gorm:"not null"`             // хеш пароля
	FullName     string `gorm:"not null"`             // ФИО
	Role         Role   `gorm:"not null"`             // роль пользователя
	CreatedAt    time.Time
}

// TeacherProfile is a derived entity that embeds User and adds the
// teacher's department. Go has no class inheritance, so embedding is the
// idiomatic equivalent — TeacherProfile gets all of User's fields plus
// its own.
type TeacherProfile struct {
	User
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

// Department is a кафедра of the educational organization.
type Department struct {
	ID   uint
	Name string `gorm:"not null"`
}

// Discipline is a учебная дисциплина tied to a department.
type Discipline struct {
	ID           uint
	Name         string `gorm:"not null"`
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

// ConsultationSlot is a time interval a teacher opens for a consultation
// on a discipline.
//
// It links a Teacher, a Discipline, and a start/end time — exactly as
// described in the domain model of the ТЗ (Приложение Б).
type ConsultationSlot struct {
	ID           uint
	TeacherID    uint `gorm:"not null;index"`
	Teacher      User `gorm:"foreignKey:TeacherID"`
	DisciplineID uint
	Discipline   Discipline `gorm:"foreignKey:DisciplineID"`
	StartsAt     time.Time  `gorm:"not null"`
	EndsAt       time.Time  `gorm:"not null"`
	Capacity     int        `gorm:"not null;default:1"` // число студентов, вмещаемых слотом
}

// BookingStatus is the status of a student's booking.
type BookingStatus string

const (
	BookingActive    BookingStatus = "active"    // запись активна
	BookingCancelled BookingStatus = "cancelled" // запись отменена студентом
	BookingDone      BookingStatus = "done"      // консультация проведена
)

// Booking links a student to a specific consultation and tracks its
// status. A Notification is created whenever that status changes.
type Booking struct {
	ID        uint
	StudentID uint             `gorm:"not null;index"`
	Student   User             `gorm:"foreignKey:StudentID"`
	SlotID    uint             `gorm:"not null;index"`
	Slot      ConsultationSlot `gorm:"foreignKey:SlotID"`
	Status    BookingStatus    `gorm:"not null;default:active"`
	CreatedAt time.Time
}

// Notification informs a user about an event tied to a booking or a
// consultation slot (created, changed, cancelled — п. 3.1.1 ТЗ).
type Notification struct {
	ID        uint
	UserID    uint   `gorm:"not null;index"`
	Message   string `gorm:"not null"`
	IsRead    bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
}
