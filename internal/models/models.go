package models

import "time"

// Role is a user role in the system.
type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

// User is the generalized entity holding login data (ТЗ Приложение Б).
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	FullName     string `gorm:"not null"`
	Role         Role   `gorm:"not null"`
	CreatedAt    time.Time
}

// Department — кафедра.
type Department struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

// Discipline — дисциплина.
type Discipline struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

// ConsultationSlot — консультационный слот, созданный преподавателем.
type ConsultationSlot struct {
	ID           uint `gorm:"primaryKey"`
	TeacherID    uint `gorm:"not null;index"`
	Teacher      User `gorm:"foreignKey:TeacherID"`
	DisciplineID uint
	Discipline   Discipline `gorm:"foreignKey:DisciplineID"`
	StartsAt     time.Time  `gorm:"not null"`
	EndsAt       time.Time  `gorm:"not null"`
	Capacity     int        `gorm:"not null;default:1"`
}

// BookingStatus — статус записи на консультацию.
type BookingStatus string

const (
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
	BookingDone      BookingStatus = "done"
)

// Booking — запись студента на консультацию.
type Booking struct {
	ID        uint             `gorm:"primaryKey"`
	StudentID uint             `gorm:"not null;index"`
	Student   User             `gorm:"foreignKey:StudentID"`
	SlotID    uint             `gorm:"not null;index"`
	Slot      ConsultationSlot `gorm:"foreignKey:SlotID"`
	Status    BookingStatus    `gorm:"not null;default:active"`
	CreatedAt time.Time
}

// Notification — уведомление пользователю о событии записи/слота.
type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Message   string `gorm:"not null"`
	IsRead    bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
}
