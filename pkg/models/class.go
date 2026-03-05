package models

import "time"

type Class struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Name         string    `json:"name" gorm:"uniqueIndex;not null"`
	Level        string    `json:"level" gorm:"type:varchar(20);not null"`
	Homeroom     string    `json:"homeroom"`
	AcademicYear string    `json:"academic_year" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateClass struct {
	Name         string `json:"name" binding:"required"`
	Level        string `json:"level" binding:"required"`
	Homeroom     string `json:"homeroom"`
	AcademicYear string `json:"academic_year" binding:"required"`
}

type UpdateClass struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	Homeroom     string `json:"homeroom"`
	AcademicYear string `json:"academic_year"`
}
