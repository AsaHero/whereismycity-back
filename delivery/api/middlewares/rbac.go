package middlewares

import (
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/entity"
	"github.com/AsaHero/whereismycity/internal/service/auth"
	"github.com/gin-gonic/gin"
)

func RoleRequired(authService auth.AuthService, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleVal, exists := c.Get("role")
		if !exists {
			outerr.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		role, ok := userRoleVal.(string)
		if !ok {
			outerr.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		hasPermission := false
		for _, r := range roles {
			if role == r || role == string(entity.UserRoleAdmin) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			outerr.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
