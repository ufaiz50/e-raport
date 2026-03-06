package models

import "time"

type SchoolProfile struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	SchoolID       *uint     `json:"school_id,omitempty" gorm:"index;unique"`
	School         *School   `json:"school,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SchoolName     string    `json:"school_name" gorm:"not null"`
	NPSN           string    `json:"npsn"`
	Address        string    `json:"address"`
	PrincipalName  string    `json:"principal_name"`
	PrincipalNIP   string    `json:"principal_nip"`
	HeadmasterSign string    `json:"headmaster_sign"`
	SchoolStamp    string    `json:"school_stamp"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UpsertSchoolProfile struct {
	SchoolID       *uint  `json:"school_id"`
	SchoolName     string `json:"school_name" binding:"required"`
	NPSN           string `json:"npsn"`
	Address        string `json:"address"`
	PrincipalName  string `json:"principal_name"`
	PrincipalNIP   string `json:"principal_nip"`
	HeadmasterSign string `json:"headmaster_sign"`
	SchoolStamp    string `json:"school_stamp"`
}
