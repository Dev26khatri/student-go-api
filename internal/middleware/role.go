package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authentication required",
			})

			c.Abort()
			return
		}
		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Success": false,
				"message": "Invalid user role",
			})
			c.Abort()
			return
		}

		for _, allowedRoles := range allowedRoles {
			if role == allowedRoles {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{
			"Sucess":  false,
			"Message": "You do not have permision to perform this action",
		})
	}
}
