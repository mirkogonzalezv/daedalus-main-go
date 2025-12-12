package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware para obtener el global_role
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("global_role")
		if role != "" {
			c.Set("global_role", role)
		}
		c.Next()
	}
}

func RequiredGlobalAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("global_role")

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
