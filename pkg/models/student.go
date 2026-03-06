package models

import "time"

type Student struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"index:idx_student_email_school,unique"`
	SchoolID  *uint     `json:"school_id,omitempty" gorm:"index:idx_student_email_school,unique"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ClassID   *uint     `json:"class_id,omitempty"`
	Class     *Class    `json:"class,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Type      string    `json:"type" gorm:"column:student_type;type:varchar(20);not null" enums:"junior,senior" example:"junior"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateStudent struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Type     string `json:"type" binding:"required,oneof=junior senior" enums:"junior,senior" example:"junior"`
	SchoolID *uint  `json:"school_id"`
	ClassID  *uint  `json:"class_id"`
}

type UpdateStudent struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Type     string `json:"type" binding:"omitempty,oneof=junior senior" enums:"junior,senior" example:"senior"`
	SchoolID *uint  `json:"school_id"`
	ClassID  *uint  `json:"class_id"`
}
