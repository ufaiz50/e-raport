package models

import "github.com/google/uuid"

type PublicUUID struct {
	UUID string `json:"uuid" gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
}

func ensureUUID(current string) string {
	if current != "" {
		return current
	}
	return uuid.NewString()
}
