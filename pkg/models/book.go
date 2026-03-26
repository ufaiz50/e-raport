package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	PublicUUID
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	SchoolID  *uint     `json:"school_id,omitempty" gorm:"index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TeacherID *uint     `json:"teacher_id,omitempty" gorm:"index"`
	Teacher   *User     `json:"teacher,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentID *uint     `json:"student_id,omitempty"`
	Student   *Student  `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (b *Book) BeforeCreate(_ *gorm.DB) error { b.UUID = ensureUUID(b.UUID); return nil }

type CreateBook struct {
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author" binding:"required"`
	SchoolID  *uint  `json:"school_id"`
	TeacherID *uint  `json:"teacher_id"`
	StudentID *uint  `json:"student_id"`
}

type UpdateBook struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	SchoolID  *uint  `json:"school_id"`
	TeacherID *uint  `json:"teacher_id"`
	StudentID *uint  `json:"student_id"`
}
