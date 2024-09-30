package messaging

import (
	"context"
	"github.com/evernetproto/evernet/internal/app/vertex/actor"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type OutboxHandler struct {
	router        *gin.Engine
	authenticator *actor.Authenticator
	manager       *OutboxManager
}

func NewOutboxHandler(router *gin.Engine, authenticator *actor.Authenticator, manager *OutboxManager) *OutboxHandler {
	return &OutboxHandler{router: router, authenticator: authenticator, manager: manager}
}

func (h *OutboxHandler) Register() {

	h.router.POST("/api/v1/messaging/outboxes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request OutboxCreationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		outbox, err := h.manager.Create(ctx, &request, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, outbox)
	})

	h.router.GET("/api/v1/messaging/outboxes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		page, size := api.Page(c)

		outboxes, err := h.manager.List(ctx, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier, page, size)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, outboxes)
	})

	h.router.GET("/api/v1/messaging/outboxes/:outboxIdentifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		outboxIdentifier := c.Param("outboxIdentifier")

		outbox, err := h.manager.Get(ctx, outboxIdentifier, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, outbox)
	})

	h.router.PUT("/api/v1/messaging/outboxes/:outboxIdentifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		outboxIdentifier := c.Param("outboxIdentifier")
		var request OutboxUpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.Update(ctx, outboxIdentifier, &request, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "outbox updated successfully")
	})
}
