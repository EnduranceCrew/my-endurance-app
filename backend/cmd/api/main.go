// cmd/api/main.go — ponto de entrada: wiring de toda a aplicação (DI manual)
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"endurance/config"
	appAlert "endurance/internal/application/alert"
	appAuth "endurance/internal/application/auth"
	appComputer "endurance/internal/application/computer"
	"endurance/internal/application/dashboard"
	appLab "endurance/internal/application/lab"
	appUser "endurance/internal/application/user"
	httpInfra "endurance/internal/infrastructure/http"
	"endurance/internal/infrastructure/http/handler"
	"endurance/internal/infrastructure/persistence"
	"endurance/internal/infrastructure/security"
)

func main() {
	// ── 1. Configuração ──────────────────────────────────────────────────────
	config.Load()
	gin.SetMode(config.App.GinMode)

	// ── 2. Banco de dados ────────────────────────────────────────────────────
	config.ConnectDB()
	persistence.Migrate(config.DB)

	// ── 3. Infra: repositórios (secondary adapters) ──────────────────────────
	userRepo     := persistence.NewUserRepository()
	labRepo      := persistence.NewLabRepository()
	computerRepo := persistence.NewComputerRepository()
	alertRepo    := persistence.NewAlertRepository()

	// ── 4. Infra: serviços de segurança ──────────────────────────────────────
	hashSvc := security.NewHashService()
	jwtSvc  := security.NewJWTService()

	// ── 5. Casos de uso (application layer) ─────────────────────────────────
	authUC      := appAuth.NewUseCase(userRepo, hashSvc, jwtSvc)
	userUC      := appUser.NewUseCase(userRepo, hashSvc)
	labUC       := appLab.NewUseCase(labRepo)
	computerUC  := appComputer.NewUseCase(computerRepo)
	alertUC     := appAlert.NewUseCase(alertRepo)
	dashboardUC := dashboard.NewUseCase(labRepo, computerRepo, userRepo, alertRepo)

	// ── 6. Seed: cria admin padrão se banco estiver vazio ───────────────────
	seedAdmin(authUC)

	// ── 7. Handlers HTTP (primary adapters) ─────────────────────────────────
	handlers := httpInfra.Handlers{
		Auth:      handler.NewAuthHandler(authUC),
		User:      handler.NewUserHandler(userUC),
		Lab:       handler.NewLabHandler(labUC),
		Computer:  handler.NewComputerHandler(computerUC),
		Alert:     handler.NewAlertHandler(alertUC),
		Dashboard: handler.NewDashboardHandler(dashboardUC),
	}

	// ── 8. Roteador ──────────────────────────────────────────────────────────
	router := httpInfra.NewRouter(handlers, jwtSvc)

	port := ":" + config.App.Port
	log.Printf("🚀 Endurance rodando em http://localhost%s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
		os.Exit(1)
	}
}

// seedAdmin cria o administrador padrão se não existir nenhum usuário.
func seedAdmin(authUC appAuth.UseCase) {
	_, err := authUC.Register(appAuth.RegisterInput{
		Name:     "Administrador",
		Email:    "admin@endurance.dev",
		CPF:      "529.982.247-25", // CPF válido para seed
		Password: "Admin@12345",
	})
	if err != nil {
		// Já existe → silencia o erro de conflito
		return
	}
	log.Println("🔑 Admin padrão criado: admin@endurance.dev / Admin@12345")
}
