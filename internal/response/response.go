package response

import (
	"github.com/gin-gonic/gin"
)

// APISuccessResponse defines the standard successful API response.
type APISuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ApiErrorResponse defines the standard failed API response.
type ApiErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Success writes a successful JSON response to the client.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APISuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes an error JSON response to the client.
func Error(c *gin.Context, status int, error string, message string) {
	c.JSON(status, ApiErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
	})
}
