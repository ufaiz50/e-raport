package models

import "time"

type Student struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email" gorm:"index:idx_student_email_school,unique"`
	NIS         string     `json:"nis"`
	NISN        string     `json:"nisn"`
	Gender      string     `json:"gender"`
	BirthPlace  string     `json:"birth_place"`
	BirthDate   *time.Time `json:"birth_date"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Religion    string     `json:"religion"`
	ParentName  string     `json:"parent_name"`
	ParentPhone string     `json:"parent_phone"`
	Status      string     `json:"status" gorm:"type:varchar(20);default:active"`
	SchoolID    *uint      `json:"school_id,omitempty" gorm:"index:idx_student_email_school,unique"`
	School      *School    `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ClassID     *uint      `json:"class_id,omitempty"`
	Class       *Class     `json:"class,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Type        string     `json:"-" gorm:"column:student_type;type:varchar(20);not null;default:junior"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateStudent struct {
	FirstName   string     `json:"first_name" binding:"required"`
	LastName    string     `json:"last_name" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	NIS         string     `json:"nis"`
	NISN        string     `json:"nisn"`
	Gender      string     `json:"gender"`
	BirthPlace  string     `json:"birth_place"`
	BirthDate   *time.Time `json:"birth_date"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Religion    string     `json:"religion"`
	ParentName  string     `json:"parent_name"`
	ParentPhone string     `json:"parent_phone"`
	Status      string     `json:"status"`
	SchoolID    *uint      `json:"school_id"`
	ClassID     *uint      `json:"class_id" binding:"required"`
}

type UpdateStudent struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email" binding:"omitempty,email"`
	NIS         string     `json:"nis"`
	NISN        string     `json:"nisn"`
	Gender      string     `json:"gender"`
	BirthPlace  string     `json:"birth_place"`
	BirthDate   *time.Time `json:"birth_date"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Religion    string     `json:"religion"`
	ParentName  string     `json:"parent_name"`
	ParentPhone string     `json:"parent_phone"`
	Status      string     `json:"status"`
	SchoolID    *uint      `json:"school_id"`
	ClassID     *uint      `json:"class_id" binding:"required"`
}
