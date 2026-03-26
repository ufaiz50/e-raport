package api

import (
	"strings"

	"golang-rest-api-template/pkg/database"
)

func whereByIDOrUUID(db database.Database, value string, schoolID *uint) database.Database {
	value = strings.TrimSpace(value)
	query := db
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	return query.Where("CAST(id AS text) = ? OR uuid = ?", value, value)
}
