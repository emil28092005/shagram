package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUsernameKey = "username"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(h, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expected Bearer token"})
			return
		}

		raw := strings.TrimSpace(strings.TrimPrefix(h, prefix))
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			return
		}

		claims, err := ParseAccessToken(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(CtxUsernameKey, claims.Username)
		c.Next()
	}
}
