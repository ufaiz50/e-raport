package models

import "time"

type ReportCardStatus string

const (
	ReportCardDraft     ReportCardStatus = "draft"
	ReportCardFinalized ReportCardStatus = "finalized"
)

type ReportCard struct {
	ID           uint             `json:"id" gorm:"primary_key"`
	StudentID    uint             `json:"student_id" gorm:"index:idx_report_term_student,unique"`
	Semester     int              `json:"semester" gorm:"index:idx_report_term_student,unique"`
	AcademicYear string           `json:"academic_year" gorm:"type:varchar(20);index:idx_report_term_student,unique"`
	Status       ReportCardStatus `json:"status" gorm:"type:varchar(20);not null;default:draft"`
	FinalizedAt  *time.Time       `json:"finalized_at"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}
