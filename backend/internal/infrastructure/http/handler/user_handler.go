package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appUser "endurance/internal/application/user"
	"endurance/pkg/response"
)

type UserHandler struct {
	useCase appUser.UseCase
}

func NewUserHandler(uc appUser.UseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	var input appUser.PaginationInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	out, err := h.useCase.GetAll(input.Page, input.Limit)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, out)
}

func (h *UserHandler) GetByID(c *gin.Context) {
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

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	var input appUser.UpdateInput
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

func (h *UserHandler) Delete(c *gin.Context) {
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

func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr, _ := c.Get("userID")
	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		response.BadRequest(c, "id inválido")
		return
	}
	var input appUser.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.useCase.ChangePassword(id, input); err != nil {
		handleError(c, err)
		return
	}
	response.Message(c, "senha alterada com sucesso")
}
