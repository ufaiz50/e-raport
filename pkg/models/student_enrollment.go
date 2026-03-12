package models

import "time"

type StudentEnrollment struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	SchoolID     *uint      `json:"school_id,omitempty" gorm:"index;not null"`
	School       *School    `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentID    uint       `json:"student_id" gorm:"index;not null"`
	Student      *Student   `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ClassID      uint       `json:"class_id" gorm:"index;not null"`
	Class        *Class     `json:"class,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AcademicYear string     `json:"academic_year" gorm:"type:varchar(20);index:idx_student_term,unique;not null"`
	Semester     int        `json:"semester" gorm:"index:idx_student_term,unique;not null"`
	IsActive     bool       `json:"is_active" gorm:"index;default:true"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateEnrollment struct {
	SchoolID     *uint      `json:"school_id"`
	StudentID    uint       `json:"student_id" binding:"required"`
	ClassID      uint       `json:"class_id" binding:"required"`
	AcademicYear string     `json:"academic_year" binding:"required"`
	Semester     int        `json:"semester" binding:"required,min=1,max=2"`
	StartDate    *time.Time `json:"start_date"`
}
