package models

import "time"

type Attendance struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	StudentID      uint      `json:"student_id" gorm:"index:idx_attendance_term_student,unique"`
	Semester       int       `json:"semester" gorm:"index:idx_attendance_term_student,unique"`
	AcademicYear   string    `json:"academic_year" gorm:"type:varchar(20);index:idx_attendance_term_student,unique"`
	SickDays       int       `json:"sick_days"`
	PermissionDays int       `json:"permission_days"`
	AbsentDays     int       `json:"absent_days"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UpsertAttendance struct {
	StudentID      uint   `json:"student_id" binding:"required"`
	Semester       int    `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear   string `json:"academic_year" binding:"required"`
	SickDays       int    `json:"sick_days" binding:"gte=0"`
	PermissionDays int    `json:"permission_days" binding:"gte=0"`
	AbsentDays     int    `json:"absent_days" binding:"gte=0"`
}
