package models

import "time"

type AcademicYear struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	SchoolID  *uint     `json:"school_id,omitempty" gorm:"index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Year      string    `json:"year" gorm:"type:varchar(20);not null"`
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Semester struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	SchoolID       *uint     `json:"school_id,omitempty" gorm:"index"`
	School         *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AcademicYearID uint      `json:"academic_year_id" gorm:"index;not null"`
	AcademicYear   *AcademicYear `json:"academic_year,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name           string    `json:"name" gorm:"type:varchar(20);not null"`
	OrderNo        int       `json:"order_no" gorm:"not null"`
	IsActive       bool      `json:"is_active" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Curriculum struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	SchoolID    *uint     `json:"school_id,omitempty" gorm:"index"`
	School      *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Year        string    `json:"year" gorm:"type:varchar(20)"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Subject struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	SchoolID  *uint     `json:"school_id,omitempty" gorm:"index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Code      string    `json:"code" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CurriculumSubject struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	CurriculumID uint       `json:"curriculum_id" gorm:"index;not null"`
	Curriculum   *Curriculum `json:"curriculum,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectID    uint       `json:"subject_id" gorm:"index;not null"`
	Subject      *Subject   `json:"subject,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type Teaching struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	SchoolID   *uint     `json:"school_id,omitempty" gorm:"index"`
	TeacherID  uint      `json:"teacher_id" gorm:"index;not null"`
	ClassID    uint      `json:"class_id" gorm:"index;not null"`
	SubjectID  uint      `json:"subject_id" gorm:"index;not null"`
	SemesterID *uint     `json:"semester_id,omitempty" gorm:"index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
