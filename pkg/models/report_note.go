package models

import "time"

type ReportNote struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	StudentID       uint      `json:"student_id" gorm:"index:idx_report_note_term_student,unique"`
	Semester        int       `json:"semester" gorm:"index:idx_report_note_term_student,unique"`
	AcademicYear    string    `json:"academic_year" gorm:"type:varchar(20);index:idx_report_note_term_student,unique"`
	HomeroomComment string    `json:"homeroom_comment" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UpsertReportNote struct {
	StudentID       uint   `json:"student_id" binding:"required"`
	Semester        int    `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear    string `json:"academic_year" binding:"required"`
	HomeroomComment string `json:"homeroom_comment" binding:"required"`
}
