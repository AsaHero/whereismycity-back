package middlewares

import (
	"strings"

	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/pkg/security"
	"github.com/gin-gonic/gin"
)

func BearerAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			outerr.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Expect the header to be "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			outerr.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Parse the JWT token.
		claims, err := security.ParseAccessToken(tokenString, secret)
		if err != nil {
			inerr.Err(err)
			outerr.Forbidden(c, "Invalid or expired token")
			c.Abort()
			return
		}

		if claims.UserID == "" {
			outerr.Forbidden(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
