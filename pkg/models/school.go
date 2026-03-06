package models

import "time"

type School struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
