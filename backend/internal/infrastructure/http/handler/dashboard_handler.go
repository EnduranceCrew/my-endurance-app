package handler

import (
	"github.com/gin-gonic/gin"
	"endurance/internal/application/dashboard"
	"endurance/pkg/response"
)

type DashboardHandler struct {
	useCase dashboard.UseCase
}

func NewDashboardHandler(uc dashboard.UseCase) *DashboardHandler {
	return &DashboardHandler{useCase: uc}
}

// GetStats retorna o resumo geral para o painel de controle.
func (h *DashboardHandler) GetStats(c *gin.Context) {
	out, err := h.useCase.GetStats()
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}
