package models

import (
	"time"

	"gorm.io/gorm"
)

type StudentEnrollment struct {
	UUIDPrimaryKey
	SchoolID     *string    `json:"school_id,omitempty" gorm:"type:uuid;index;not null"`
	School       *School    `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	StudentID    string     `json:"student_id" gorm:"type:uuid;index:idx_student_term,unique;not null"`
	Student      *Student   `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
	ClassID      string     `json:"class_id" gorm:"type:uuid;index:idx_student_term,unique;not null"`
	Class        *Class     `json:"class,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ClassID;references:ID"`
	AcademicYear string     `json:"academic_year" gorm:"type:varchar(20);index:idx_student_term,unique;not null"`
	Semester     int        `json:"semester" gorm:"index:idx_student_term,unique;not null"`
	IsActive     bool       `json:"is_active" gorm:"index;default:true"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (e *StudentEnrollment) BeforeCreate(_ *gorm.DB) error { e.ID = ensureUUID(e.ID); return nil }

type CreateEnrollment struct {
	SchoolID     *string    `json:"school_id"`
	StudentID    string     `json:"student_id" binding:"required"`
	ClassID      string     `json:"class_id" binding:"required"`
	AcademicYear string     `json:"academic_year" binding:"required"`
	Semester     int        `json:"semester" binding:"required,min=1,max=2"`
	StartDate    *time.Time `json:"start_date"`
}
