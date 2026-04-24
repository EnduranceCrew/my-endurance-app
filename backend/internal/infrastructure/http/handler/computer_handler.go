package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appComputer "endurance/internal/application/computer"
	"endurance/pkg/response"
)

type ComputerHandler struct {
	useCase appComputer.UseCase
}

func NewComputerHandler(uc appComputer.UseCase) *ComputerHandler {
	return &ComputerHandler{useCase: uc}
}

func (h *ComputerHandler) GetAll(c *gin.Context) {
	page, limit := paginationParams(c)
	statusFilter := c.Query("status")
	out, err := h.useCase.GetAll(c.Request.Context(), page, limit, statusFilter)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *ComputerHandler) GetByLabID(c *gin.Context) {
	labID, err := uuid.Parse(c.Param("labId"))
	if err != nil {
		response.BadRequest(c, "lab_id inválido")
		return
	}
	out, err := h.useCase.GetByLabID(c.Request.Context(), labID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *ComputerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	out, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *ComputerHandler) Create(c *gin.Context) {
	var input appComputer.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.useCase.Create(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, out)
}

func (h *ComputerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	var input appComputer.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.useCase.Update(c.Request.Context(), id, input)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *ComputerHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	var input appComputer.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.useCase.UpdateStatus(c.Request.Context(), id, input); err != nil {
		handleError(c, err)
		return
	}
	response.Message(c, "status atualizado")
}

func (h *ComputerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	response.NoContent(c)
}
