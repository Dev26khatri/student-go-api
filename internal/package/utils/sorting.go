package utils

import "strings"

// GetSortColumn returns a safe database column for the requested sort field.
func GetSortColumn(sortBy string) string {
	allowedColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"age":        "age",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	if column, exists := allowedColumns[sortBy]; exists {
		return column
	}
	return "id"
}

// GetSortOrder returns a safe ascending or descending sort direction.
func GetSortOrder(order string) string {
	order = strings.ToLower(strings.TrimSpace(order))
	if order == "asc" {
		return "ASC"
	}
	return "DESC"
}
