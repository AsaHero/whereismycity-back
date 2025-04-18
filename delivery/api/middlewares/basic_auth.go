package middlewares

import (
	"github.com/AsaHero/whereismycity/delivery/api/outerr"
	"github.com/AsaHero/whereismycity/internal/service/auth"
	"github.com/gin-gonic/gin"
)

func BasicAuth(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			c.Header("WWW-Authenticate", "Basic realm=Restricted")
			outerr.Unauthorized(c, "Unauthorized")

			c.Abort()
			return
		}

		user, err := authService.LoginByUsername(c, username, password)
		if err != nil {
			outerr.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("role", string(user.Role))

		c.Next()
	}
}
