package models

import (
	"time"

	"gorm.io/gorm"
)

type Grade struct {
	UUIDPrimaryKey
	SchoolID       *string            `json:"school_id,omitempty" gorm:"type:uuid;index"`
	School         *School            `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:SchoolID;references:ID"`
	SemesterID     *string            `json:"semester_id,omitempty" gorm:"type:uuid;index"`
	SemesterRef    *Semester          `json:"semester_ref,omitempty" gorm:"foreignKey:SemesterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TeachingID     *string            `json:"teaching_id,omitempty" gorm:"type:uuid;index"`
	Teaching       *Teaching          `json:"teaching,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:TeachingID;references:ID"`
	EnrollmentID   *string            `json:"enrollment_id,omitempty" gorm:"type:uuid;index"`
	Enrollment     *StudentEnrollment `json:"enrollment,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:EnrollmentID;references:ID"`
	StudentID      string             `json:"student_id" gorm:"type:uuid"`
	Student        Student            `json:"student" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
	BookID         string             `json:"book_id" gorm:"type:uuid"`
	Book           Book               `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:BookID;references:ID"`
	Semester       int                `json:"semester" gorm:"not null"`
	AcademicYear   string             `json:"academic_year" gorm:"type:varchar(20);not null"`
	KnowledgeScore float64            `json:"knowledge_score" gorm:"type:numeric(5,2);not null"`
	SkillScore     float64            `json:"skill_score" gorm:"type:numeric(5,2);not null"`
	FinalScore     float64            `json:"final_score" gorm:"type:numeric(5,2);not null"`
	Notes          string             `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time          `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time          `json:"updated_at" gorm:"autoUpdateTime"`
}

func (g *Grade) BeforeCreate(_ *gorm.DB) error { g.ID = ensureUUID(g.ID); return nil }

type CreateGrade struct {
	SchoolID       *string  `json:"school_id"`
	SemesterID     *string  `json:"semester_id"`
	TeachingID     *string  `json:"teaching_id"`
	EnrollmentID   *string  `json:"enrollment_id"`
	StudentID      string   `json:"student_id" binding:"required"`
	BookID         string   `json:"book_id" binding:"required"`
	Semester       int      `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear   string   `json:"academic_year" binding:"required"`
	KnowledgeScore float64  `json:"knowledge_score" binding:"required,gte=0,lte=100"`
	SkillScore     float64  `json:"skill_score" binding:"required,gte=0,lte=100"`
	Notes          string   `json:"notes"`
}

type UpdateGrade struct {
	SchoolID       *string  `json:"school_id"`
	SemesterID     *string  `json:"semester_id"`
	TeachingID     *string  `json:"teaching_id"`
	EnrollmentID   *string  `json:"enrollment_id"`
	KnowledgeScore *float64 `json:"knowledge_score" binding:"omitempty,gte=0,lte=100"`
	SkillScore     *float64 `json:"skill_score" binding:"omitempty,gte=0,lte=100"`
	Notes          *string  `json:"notes"`
}
