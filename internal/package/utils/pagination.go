package utils

// Pagination stores normalized page, limit, and offset values.
type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

// PaginationMeta describes the result count and available pages.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_page"`
}

// NewPagination normalizes paging values and calculates the database offset.
func NewPagination(page, limit int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// NewPaginationMeta calculates pagination metadata from the total result count.
func NewPaginationMeta(p Pagination, total int64) PaginationMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
	}
	return PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
