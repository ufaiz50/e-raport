package models

import (
	"time"

	"gorm.io/gorm"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterUser struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Role     string  `json:"role" binding:"omitempty,oneof=super_admin admin guru wali_kelas"`
	SchoolID *string `json:"school_id"`
}

type User struct {
	UUIDPrimaryKey
	Username  string    `json:"username" gorm:"index:idx_username_school,unique"`
	SchoolID  *string   `json:"school_id,omitempty" gorm:"type:uuid;index:idx_username_school,unique"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	Password  string    `json:"password"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null;default:guru"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error { u.ID = ensureUUID(u.ID); return nil }
