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
}
