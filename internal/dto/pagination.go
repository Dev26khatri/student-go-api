package dto

// paginationMeta describes pagination information returned by an API.
type paginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_page"`
}
