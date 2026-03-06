package models

import "time"

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	SchoolID *uint  `json:"school_id"`
}

type RegisterUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"omitempty,oneof=super_admin admin guru wali_kelas"`
	SchoolID *uint  `json:"school_id"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"index:idx_username_school,unique"`
	SchoolID  *uint     `json:"school_id,omitempty" gorm:"index:idx_username_school,unique"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Password  string    `json:"password"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null;default:guru"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
