package models

import (
	"time"

	"gorm.io/gorm"
)

type AcademicYear struct {
	UUIDPrimaryKey
	SchoolID  *string   `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	Year      string    `json:"year" gorm:"type:varchar(20);not null"`
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Semester struct {
	UUIDPrimaryKey
	SchoolID       *string       `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School         *School       `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	AcademicYearID string        `json:"academic_year_id" gorm:"type:uuid;index;not null"`
	AcademicYear   *AcademicYear `json:"academic_year,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AcademicYearID;references:ID"`
	Name           string        `json:"name" gorm:"type:varchar(20);not null"`
	OrderNo        int           `json:"order_no" gorm:"not null"`
	IsActive       bool          `json:"is_active" gorm:"default:false"`
	CreatedAt      time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type Curriculum struct {
	UUIDPrimaryKey
	SchoolID    *string   `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School      *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Year        string    `json:"year" gorm:"type:varchar(20)"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Subject struct {
	UUIDPrimaryKey
	SchoolID  *string   `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School    *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Code      string    `json:"code" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CurriculumSubject struct {
	UUIDPrimaryKey
	CurriculumID string      `json:"curriculum_id" gorm:"type:uuid;index;not null"`
	Curriculum   *Curriculum `json:"curriculum,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:CurriculumID;references:ID"`
	SubjectID    string      `json:"subject_id" gorm:"type:uuid;index;not null"`
	Subject      *Subject    `json:"subject,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SubjectID;references:ID"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type Teaching struct {
	UUIDPrimaryKey
	SchoolID   *string   `json:"school_id,omitempty" gorm:"type:uuid;index"`
	TeacherID  string    `json:"teacher_id" gorm:"type:uuid;index;not null"`
	ClassID    string    `json:"class_id" gorm:"type:uuid;index;not null"`
	SubjectID  string    `json:"subject_id" gorm:"type:uuid;index;not null"`
	SemesterID *string   `json:"semester_id,omitempty" gorm:"type:uuid;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (m *AcademicYear) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
func (m *Semester) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
func (m *Curriculum) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
func (m *Subject) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
func (m *CurriculumSubject) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
func (m *Teaching) BeforeCreate(_ *gorm.DB) error { m.ID = ensureUUID(m.ID); return nil }
