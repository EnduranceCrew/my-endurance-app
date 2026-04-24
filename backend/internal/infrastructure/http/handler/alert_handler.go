package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appAlert "endurance/internal/application/alert"
	"endurance/pkg/response"
)

type AlertHandler struct {
	useCase appAlert.UseCase
}

func NewAlertHandler(uc appAlert.UseCase) *AlertHandler {
	return &AlertHandler{useCase: uc}
}

func (h *AlertHandler) GetAll(c *gin.Context) {
	onlyOpen := c.DefaultQuery("open", "true") == "true"
	page, limit := paginationParams(c)
	out, err := h.useCase.GetAll(onlyOpen, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *AlertHandler) GetByLabID(c *gin.Context) {
	labID, err := uuid.Parse(c.Param("labId"))
	if err != nil {
		response.BadRequest(c, "lab_id inválido")
		return
	}
	onlyOpen := c.DefaultQuery("open", "true") == "true"
	out, err := h.useCase.GetByLabID(labID, onlyOpen)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *AlertHandler) Create(c *gin.Context) {
	var input appAlert.CreateInput
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

func (h *AlertHandler) Resolve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	if err := h.useCase.Resolve(id); err != nil {
		handleError(c, err)
		return
	}
	response.Message(c, "alerta resolvido")
}

func (h *AlertHandler) Delete(c *gin.Context) {
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

func (h *AlertHandler) BulkResolve(c *gin.Context) {
	var input appAlert.BulkResolveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	count, err := h.useCase.BulkResolve(input)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, gin.H{"resolved": count})
}
