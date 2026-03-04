package models

import "time"

type Book struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	StudentID *uint     `json:"student_id,omitempty"`
	Student   *Student  `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateBook struct {
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author" binding:"required"`
	StudentID *uint  `json:"student_id"`
}

type UpdateBook struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	StudentID *uint  `json:"student_id"`
}
