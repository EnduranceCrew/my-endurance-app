package middleware

import (
	"github.com/gin-gonic/gin"
	"endurance/pkg/response"
)

// RequireRole retorna 403 se o usuário autenticado não tiver o role exigido.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			response.Unauthorized(c, "não autenticado")
			return
		}

		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "permissão insuficiente")
	}
}

// AdminOnly é um alias conveniente para RequireRole("admin").
func AdminOnly() gin.HandlerFunc {
	return RequireRole("admin")
}
