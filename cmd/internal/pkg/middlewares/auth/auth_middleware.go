package middlewares

import (
	"net/http"
	"strings"

	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware para obtener el global_role
func AuthMiddleware(authService *authService.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Validamos header Authorization
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Confirmamos el formato Bearer
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Validamos token+
		claims, err := authService.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func RequiredGlobalAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "root" && role != "system_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "System admin role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequiredTenantAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "root" && role != "admin" && role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tenant admin role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
