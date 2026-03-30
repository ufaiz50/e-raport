package models

import "github.com/google/uuid"

type UUIDPrimaryKey struct {
	ID string `json:"id" gorm:"type:uuid;primaryKey"`
}

func ensureUUID(current string) string {
	if current != "" {
		return current
	}
	return uuid.NewString()
}
