// cmd/api/main.go — ponto de entrada: wiring de toda a aplicação (DI manual)
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"

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

	if config.App.JWTSecret == "secret_dev_only" {
		slog.Warn("⚠️  JWT_SECRET está usando valor padrão inseguro! Defina um segredo forte em produção.")
	}

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
	slog.Info("🚀 Endurance rodando", "addr", "http://localhost"+port)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("falha ao iniciar servidor", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("desligando servidor graciosamente...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("erro ao desligar servidor", "error", err)
	}
	slog.Info("servidor encerrado")
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
