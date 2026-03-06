package models

import "time"

type ReportCardStatus string

const (
	ReportCardDraft     ReportCardStatus = "draft"
	ReportCardFinalized ReportCardStatus = "finalized"
)

type ReportCard struct {
	ID           uint             `json:"id" gorm:"primary_key"`
	SchoolID     *uint            `json:"school_id,omitempty" gorm:"index:idx_report_term_student_school,unique"`
	School       *School          `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentID    uint             `json:"student_id" gorm:"index:idx_report_term_student_school,unique"`
	Semester     int              `json:"semester" gorm:"index:idx_report_term_student_school,unique"`
	AcademicYear string           `json:"academic_year" gorm:"type:varchar(20);index:idx_report_term_student_school,unique"`
	Status       ReportCardStatus `json:"status" gorm:"type:varchar(20);not null;default:draft"`
	FinalizedAt  *time.Time       `json:"finalized_at"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}
