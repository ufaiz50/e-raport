package models

import "time"

type Grade struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	SchoolID       *uint     `json:"school_id,omitempty" gorm:"index"`
	School         *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentID      uint      `json:"student_id"`
	Student        Student   `json:"student" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BookID         uint      `json:"book_id"`
	Book           Book      `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Semester       int       `json:"semester" gorm:"not null"`
	AcademicYear   string    `json:"academic_year" gorm:"type:varchar(20);not null"`
	KnowledgeScore float64   `json:"knowledge_score" gorm:"type:numeric(5,2);not null"`
	SkillScore     float64   `json:"skill_score" gorm:"type:numeric(5,2);not null"`
	FinalScore     float64   `json:"final_score" gorm:"type:numeric(5,2);not null"`
	Notes          string    `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateGrade struct {
	SchoolID       *uint   `json:"school_id"`
	StudentID      uint    `json:"student_id" binding:"required"`
	BookID         uint    `json:"book_id" binding:"required"`
	Semester       int     `json:"semester" binding:"required,min=1,max=2"`
	AcademicYear   string  `json:"academic_year" binding:"required"`
	KnowledgeScore float64 `json:"knowledge_score" binding:"required,gte=0,lte=100"`
	SkillScore     float64 `json:"skill_score" binding:"required,gte=0,lte=100"`
	Notes          string  `json:"notes"`
}

type UpdateGrade struct {
	SchoolID       *uint    `json:"school_id"`
	KnowledgeScore *float64 `json:"knowledge_score" binding:"omitempty,gte=0,lte=100"`
	SkillScore     *float64 `json:"skill_score" binding:"omitempty,gte=0,lte=100"`
	Notes          *string  `json:"notes"`
}
