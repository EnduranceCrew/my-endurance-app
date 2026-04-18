package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appAuth "endurance/internal/application/auth"
	"endurance/pkg/apperrors"
	"endurance/pkg/response"
)

type AuthHandler struct {
	useCase appAuth.UseCase
}

func NewAuthHandler(uc appAuth.UseCase) *AuthHandler {
	return &AuthHandler{useCase: uc}
}

// Login godoc
// @Summary      Autenticar usuário
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body appAuth.LoginInput true "Credenciais"
// @Success      200  {object}  appAuth.TokenOutput
// @Failure      401  {object}  response.Envelope
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input appAuth.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	out, err := h.useCase.Login(input)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, out)
}

// Register godoc
// @Summary      Criar conta de usuário
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body appAuth.RegisterInput true "Dados do usuário"
// @Success      201  {object}  appAuth.TokenOutput
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input appAuth.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	out, err := h.useCase.Register(input)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Created(c, out)
}

// ── helper compartilhado entre handlers ─────────────────────────────────────

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPCode, appErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "erro interno do servidor")
}
