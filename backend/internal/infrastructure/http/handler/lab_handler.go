package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appLab "endurance/internal/application/lab"
	"endurance/pkg/response"
)

type LabHandler struct {
	useCase appLab.UseCase
}

func NewLabHandler(uc appLab.UseCase) *LabHandler {
	return &LabHandler{useCase: uc}
}

func (h *LabHandler) GetAll(c *gin.Context) {
	page, limit := paginationParams(c)
	out, err := h.useCase.GetAll(page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *LabHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	out, err := h.useCase.GetByID(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *LabHandler) Create(c *gin.Context) {
	var input appLab.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.useCase.Create(input)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, out)
}

func (h *LabHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	var input appLab.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.useCase.Update(id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *LabHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	if err := h.useCase.Delete(id); err != nil {
		handleError(c, err)
		return
	}
	response.NoContent(c)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func paginationParams(c *gin.Context) (int, int) {
	page := 1
	limit := 20
	if p := c.DefaultQuery("page", "1"); p != "" {
		if v, err := parseInt(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.DefaultQuery("limit", "20"); l != "" {
		if v, err := parseInt(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	return page, limit
}

func parseInt(s string) (int, error) {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v) //nolint
	return v, err
}
