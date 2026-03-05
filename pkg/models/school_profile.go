package models

import "time"

type SchoolProfile struct {
	ID             uint      `json:"id" gorm:"primary_key"`
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
	SchoolName     string `json:"school_name" binding:"required"`
	NPSN           string `json:"npsn"`
	Address        string `json:"address"`
	PrincipalName  string `json:"principal_name"`
	PrincipalNIP   string `json:"principal_nip"`
	HeadmasterSign string `json:"headmaster_sign"`
	SchoolStamp    string `json:"school_stamp"`
}
