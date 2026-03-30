package models

import (
	"time"

	"gorm.io/gorm"
)

type ReportNote struct {
	UUIDPrimaryKey
	SchoolID        *string            `json:"school_id,omitempty" gorm:"type:uuid;index:idx_report_note_term_student_school,unique"`
	School          *School            `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	EnrollmentID    *string            `json:"enrollment_id,omitempty" gorm:"type:uuid;index"`
	Enrollment      *StudentEnrollment `json:"enrollment,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:EnrollmentID;references:ID"`
	StudentID       string             `json:"student_id" gorm:"type:uuid;index:idx_report_note_term_student_school,unique"`
	Semester        int                `json:"semester" gorm:"index:idx_report_note_term_student_school,unique"`
	AcademicYear    string             `json:"academic_year" gorm:"type:varchar(20);index:idx_report_note_term_student_school,unique"`
	HomeroomComment string             `json:"homeroom_comment" gorm:"type:text"`
	CreatedAt       time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

func (r *ReportNote) BeforeCreate(_ *gorm.DB) error { r.ID = ensureUUID(r.ID); return nil }

type UpsertReportNote struct {
	SchoolID        *string `json:"school_id"`
	EnrollmentID    *string `json:"enrollment_id"`
	StudentID       string  `json:"student_id" binding:"required"`
	Semester        int     `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear    string  `json:"academic_year" binding:"required"`
	HomeroomComment string  `json:"homeroom_comment" binding:"required"`
}
