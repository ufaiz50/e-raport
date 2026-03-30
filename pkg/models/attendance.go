package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	UUIDPrimaryKey
	SchoolID       *string            `json:"school_id,omitempty" gorm:"type:uuid;index:idx_attendance_term_student_school,unique"`
	School         *School            `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	EnrollmentID   *string            `json:"enrollment_id,omitempty" gorm:"type:uuid;index"`
	Enrollment     *StudentEnrollment `json:"enrollment,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:EnrollmentID;references:ID"`
	StudentID      string             `json:"student_id" gorm:"type:uuid;index:idx_attendance_term_student_school,unique"`
	Semester       int                `json:"semester" gorm:"index:idx_attendance_term_student_school,unique"`
	AcademicYear   string             `json:"academic_year" gorm:"type:varchar(20);index:idx_attendance_term_student_school,unique"`
	SickDays       int                `json:"sick_days"`
	PermissionDays int                `json:"permission_days"`
	AbsentDays     int                `json:"absent_days"`
	CreatedAt      time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a *Attendance) BeforeCreate(_ *gorm.DB) error { a.ID = ensureUUID(a.ID); return nil }

type UpsertAttendance struct {
	SchoolID       *string `json:"school_id"`
	EnrollmentID   *string `json:"enrollment_id"`
	StudentID      string  `json:"student_id" binding:"required"`
	Semester       int     `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear   string  `json:"academic_year" binding:"required"`
	SickDays       int     `json:"sick_days" binding:"gte=0"`
	PermissionDays int     `json:"permission_days" binding:"gte=0"`
	AbsentDays     int     `json:"absent_days" binding:"gte=0"`
}
