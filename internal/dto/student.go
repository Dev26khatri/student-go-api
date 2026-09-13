package dto

import (
	"mime/multipart"
	"time"
)

// CreateStudentRequest contains data required to create a student.
type CreateStudentRequest struct {
	Name  string                `form:"name" binding:"required"`
	Email string                `form:"email" binding:"required,email"`
	Age   int                   `form:"age" binding:"required,min=18"`
	Photo *multipart.FileHeader `form:"photo"`
}

// UpdateStudentRequest contains data used to update a student.
type UpdateStudentRequest struct {
	Name  string                `form:"name" binding:"required"`
	Email string                `form:"email" binding:"required,email"`
	Age   int                   `form:"age" binding:"required,min=18"`
	Photo *multipart.FileHeader `form:"photo"`
}

// StudentResponse defines the public student response format.
type StudentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	PhotoUrl  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationResponse combines pagination details with response data.
type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}
