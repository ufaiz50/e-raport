package api

import (
	"strings"

	"golang-rest-api-template/pkg/database"
)

func whereByIDOrUUID(db database.Database, value string, schoolID *string) database.Database {
	value = strings.TrimSpace(value)
	query := db
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	return query.Where("id = ?", value)
}
