package response

import "github.com/gin-gonic/gin"

// ── Envelope padrão de resposta ──────────────────────────────────────────────

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ── Respostas de sucesso ─────────────────────────────────────────────────────

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Envelope{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Envelope{Success: true, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(204)
}

func Message(c *gin.Context, msg string) {
	c.JSON(200, Envelope{Success: true, Message: msg})
}

// ── Respostas de erro ────────────────────────────────────────────────────────

func Error(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Envelope{Success: false, Error: message})
}

func BadRequest(c *gin.Context, msg string) { Error(c, 400, msg) }

func Unauthorized(c *gin.Context, msg string) { Error(c, 401, msg) }

func Forbidden(c *gin.Context, msg string) { Error(c, 403, msg) }

func NotFound(c *gin.Context, msg string) { Error(c, 404, msg) }

func Conflict(c *gin.Context, msg string) { Error(c, 409, msg) }

func UnprocessableEntity(c *gin.Context, msg string) { Error(c, 422, msg) }

func Internal(c *gin.Context, msg string) { Error(c, 500, msg) }
