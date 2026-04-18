// Package http monta o roteador Gin com todos os grupos de rotas.
package http

import (
	"github.com/gin-gonic/gin"

	"endurance/internal/infrastructure/http/handler"
	"endurance/internal/infrastructure/http/middleware"
	"endurance/internal/infrastructure/security"
)

// Handlers agrupa todos os handlers da aplicação.
type Handlers struct {
	Auth      *handler.AuthHandler
	User      *handler.UserHandler
	Lab       *handler.LabHandler
	Computer  *handler.ComputerHandler
	Alert     *handler.AlertHandler
	Dashboard *handler.DashboardHandler
}

// NewRouter constrói e retorna o motor Gin com todas as rotas registradas.
func NewRouter(h Handlers, jwtSvc *security.JWTServiceImpl) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORS())

	// ── Health check ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "endurance"})
	})

	// ── API v1 ────────────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// Rotas públicas
	auth := v1.Group("/auth")
	{
		auth.POST("/login", h.Auth.Login)
		auth.POST("/register", h.Auth.Register)
	}

	// Rotas protegidas (requer JWT)
	protected := v1.Group("/")
	protected.Use(middleware.Auth(jwtSvc))
	{
		// Dashboard — admin e user visualizam
		protected.GET("dashboard/stats", h.Dashboard.GetStats)

		// Usuários — apenas admin gerencia
		users := protected.Group("users")
		{
			users.GET("", middleware.AdminOnly(), h.User.GetAll)
			users.GET(":id", middleware.AdminOnly(), h.User.GetByID)
			users.PUT(":id", middleware.AdminOnly(), h.User.Update)
			users.DELETE(":id", middleware.AdminOnly(), h.User.Delete)
			users.POST("me/password", h.User.ChangePassword)
		}

		// Laboratórios
		labs := protected.Group("labs")
		{
			labs.GET("", h.Lab.GetAll)
			labs.GET(":id", h.Lab.GetByID)
			labs.POST("", middleware.AdminOnly(), h.Lab.Create)
			labs.PUT(":id", middleware.AdminOnly(), h.Lab.Update)
			labs.DELETE(":id", middleware.AdminOnly(), h.Lab.Delete)

			// Alertas por laboratório
			labs.GET(":labId/alerts", h.Alert.GetByLabID)
			// Computadores por laboratório
			labs.GET(":labId/computers", h.Computer.GetByLabID)
		}

		// Computadores
		computers := protected.Group("computers")
		{
			computers.GET("", h.Computer.GetAll)
			computers.GET(":id", h.Computer.GetByID)
			computers.POST("", middleware.AdminOnly(), h.Computer.Create)
			computers.PUT(":id", middleware.AdminOnly(), h.Computer.Update)
			computers.PATCH(":id/status", h.Computer.UpdateStatus)
			computers.DELETE(":id", middleware.AdminOnly(), h.Computer.Delete)
		}

		// Alertas
		alerts := protected.Group("alerts")
		{
			alerts.GET("", h.Alert.GetAll)
			alerts.POST("", h.Alert.Create)
			alerts.PATCH(":id/resolve", middleware.AdminOnly(), h.Alert.Resolve)
			alerts.DELETE(":id", middleware.AdminOnly(), h.Alert.Delete)
		}
	}

	return r
}
