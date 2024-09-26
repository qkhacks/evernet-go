package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	router *gin.Engine
}

func NewHandler(router *gin.Engine) *Handler {
	return &Handler{router: router}
}

func (h *Handler) Register() {
	h.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, &CheckResponse{Status: "ok"})
	})
}
