package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	UUIDPrimaryKey
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	SchoolID  *string   `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	TeacherID *string   `json:"teacher_id,omitempty" gorm:"type:uuid;index"`
	Teacher   *User     `json:"teacher,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:TeacherID;references:ID"`
	StudentID *string   `json:"student_id,omitempty" gorm:"type:uuid"`
	Student   *Student  `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:StudentID;references:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (b *Book) BeforeCreate(_ *gorm.DB) error { b.ID = ensureUUID(b.ID); return nil }

type CreateBook struct {
	Title     string  `json:"title" binding:"required"`
	Author    string  `json:"author" binding:"required"`
	SchoolID  *string `json:"school_id"`
	TeacherID *string `json:"teacher_id"`
	StudentID *string `json:"student_id"`
}

type UpdateBook struct {
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	SchoolID  *string `json:"school_id"`
	TeacherID *string `json:"teacher_id"`
	StudentID *string `json:"student_id"`
}
