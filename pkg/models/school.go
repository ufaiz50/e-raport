package models

import "time"

type School struct {
	ID             uint      `json:"id" gorm:"primary_key"`
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
