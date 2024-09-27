package admin

import (
	"context"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	router  *gin.Engine
	manager *Manager
}

func NewHandler(router *gin.Engine, manager *Manager) *Handler {
	return &Handler{router: router, manager: manager}
}

func (h *Handler) Register() {
	h.router.POST("/api/v1/admins/init", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var request InitRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		admin, err := h.manager.Init(ctx, &request)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, admin)
	})

	h.router.POST("/api/v1/admins/token", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var request TokenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		token, err := h.manager.GetToken(ctx, &request)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, token)
	})
}
