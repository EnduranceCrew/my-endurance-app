package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"endurance/internal/infrastructure/security"
	"endurance/pkg/response"
)

const UserIDKey = "userID"
const UserRoleKey = "userRole"

// Auth valida o token JWT em cada requisição protegida.
func Auth(jwtSvc *security.JWTServiceImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c, "token de autenticação não fornecido")
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtSvc.Validate(tokenStr)
		if err != nil {
			response.Unauthorized(c, "token inválido ou expirado")
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)
		c.Next()
	}
}
