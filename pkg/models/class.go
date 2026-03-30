package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	UUIDPrimaryKey
	Name         string    `json:"name" gorm:"index:idx_class_name_school,unique;not null"`
	SchoolID     *string   `json:"school_id,omitempty" gorm:"type:uuid;index:idx_class_name_school,unique"`
	School       *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	Level        string    `json:"level" gorm:"type:varchar(20);not null"`
	Homeroom     string    `json:"homeroom"`
	AcademicYear string    `json:"academic_year" gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (c *Class) BeforeCreate(_ *gorm.DB) error { c.ID = ensureUUID(c.ID); return nil }

type CreateClass struct {
	Name         string  `json:"name" binding:"required"`
	Level        string  `json:"level" binding:"required"`
	Homeroom     string  `json:"homeroom"`
	AcademicYear string  `json:"academic_year" binding:"required"`
	SchoolID     *string `json:"school_id"`
}

type UpdateClass struct {
	Name         string  `json:"name"`
	Level        string  `json:"level"`
	Homeroom     string  `json:"homeroom"`
	AcademicYear string  `json:"academic_year"`
	SchoolID     *string `json:"school_id"`
}
