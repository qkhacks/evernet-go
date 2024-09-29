package messaging

import (
	"context"
	"github.com/evernetproto/evernet/internal/app/vertex/actor"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type InboxHandler struct {
	router        *gin.Engine
	authenticator *actor.Authenticator
	manager       *InboxManager
}

func NewInboxHandler(router *gin.Engine, authenticator *actor.Authenticator, manager *InboxManager) *InboxHandler {
	return &InboxHandler{router: router, authenticator: authenticator, manager: manager}
}

func (h *InboxHandler) Register() {

	h.router.POST("/api/v1/messaging/inboxes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request InboxCreationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		inbox, err := h.manager.Create(ctx, &request, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, inbox)
	})

	h.router.GET("/api/v1/messaging/inboxes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		page, size := api.Page(c)

		inboxes, err := h.manager.List(ctx, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier, page, size)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, inboxes)
	})

	h.router.GET("/api/v1/messaging/inboxes/:inboxIdentifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("inboxIdentifier")

		inbox, err := h.manager.Get(ctx, identifier, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, inbox)
	})

	h.router.PUT("/api/v1/messaging/inboxes/:inboxIdentifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("inboxIdentifier")
		var request InboxUpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.Update(ctx, identifier, &request, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "inbox updated successfully")
	})

	h.router.DELETE("/api/v1/messaging/inboxes/:inboxIdentifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("inboxIdentifier")

		err = h.manager.Delete(ctx, identifier, authenticatedActor.Address, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "inbox deleted successfully")
	})
}
