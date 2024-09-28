package actor

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
	h.router.POST("/api/v1/nodes/:nodeIdentifier/actors/signup", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var request SignUpRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		nodeIdentifier := c.Param("nodeIdentifier")

		actor, err := h.manager.SignUp(ctx, nodeIdentifier, &request)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, actor)
	})
}
