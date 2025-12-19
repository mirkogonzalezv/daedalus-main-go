package middlewares

import (
	rateLimitService "daedalus-engine-go/cmd/internal/features/ratelimit/infraestructure/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func IPRateLimitMiddleware(rateLimitSvc *rateLimitService.RateLimitService, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		endpoint := c.Request.URL.Path
		userAgent := c.GetHeader("User-Agent")
		result, err := rateLimitSvc.CheckRateLimitByIP(c.Request.Context(), ip, endpoint, userAgent)

		if err != nil {
			log.Error("IP Rate limit falló en confirmación", zap.Error(err))
			c.Next()
			return
		}

		c.Header("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

		if !result.Allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": int64(result.RetryAfter.Seconds()),
				"message":     "Too many request from this IP address",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
