package models

import (
	"time"

	"gorm.io/gorm"
)

type School struct {
	UUIDPrimaryKey
	Name           string    `json:"name" gorm:"uniqueIndex;not null"`
	Code           string    `json:"code" gorm:"uniqueIndex;not null"`
	Address        string    `json:"address"`
	NPSN           string    `json:"npsn"`
	PrincipalName  string    `json:"principal_name"`
	PrincipalNIP   string    `json:"principal_nip"`
	HeadmasterSign string    `json:"headmaster_sign"`
	SchoolStamp    string    `json:"school_stamp"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *School) BeforeCreate(_ *gorm.DB) error { s.ID = ensureUUID(s.ID); return nil }
