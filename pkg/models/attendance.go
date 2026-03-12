package models

import "time"

type Attendance struct {
	ID             uint               `json:"id" gorm:"primary_key"`
	SchoolID       *uint              `json:"school_id,omitempty" gorm:"index:idx_attendance_term_student_school,unique"`
	School         *School            `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	EnrollmentID   *uint              `json:"enrollment_id,omitempty" gorm:"index"`
	Enrollment     *StudentEnrollment `json:"enrollment,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentID      uint               `json:"student_id" gorm:"index:idx_attendance_term_student_school,unique"`
	Semester       int                `json:"semester" gorm:"index:idx_attendance_term_student_school,unique"`
	AcademicYear   string             `json:"academic_year" gorm:"type:varchar(20);index:idx_attendance_term_student_school,unique"`
	SickDays       int                `json:"sick_days"`
	PermissionDays int                `json:"permission_days"`
	AbsentDays     int                `json:"absent_days"`
	CreatedAt      time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

type UpsertAttendance struct {
	SchoolID       *uint  `json:"school_id"`
	EnrollmentID   *uint  `json:"enrollment_id"`
	StudentID      uint   `json:"student_id" binding:"required"`
	Semester       int    `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear   string `json:"academic_year" binding:"required"`
	SickDays       int    `json:"sick_days" binding:"gte=0"`
	PermissionDays int    `json:"permission_days" binding:"gte=0"`
	AbsentDays     int    `json:"absent_days" binding:"gte=0"`
}
