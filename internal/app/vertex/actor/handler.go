package actor

import (
	"context"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	router        *gin.Engine
	authenticator *Authenticator
	manager       *Manager
}

func NewHandler(router *gin.Engine, authenticator *Authenticator, manager *Manager) *Handler {
	return &Handler{
		router:        router,
		authenticator: authenticator,
		manager:       manager,
	}
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

	h.router.POST("/api/v1/nodes/:nodeIdentifier/actors/token", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var request TokenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		nodeIdentifier := c.Param("nodeIdentifier")

		token, err := h.manager.GetToken(ctx, nodeIdentifier, &request)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, token)
	})

	h.router.GET("/api/v1/actors/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)

		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		if !authenticatedActor.IsLocal {
			api.ErrorMessage(c, http.StatusForbidden, "not allowed")
			return
		}

		actor, err := h.manager.Get(ctx, authenticatedActor.Identifier, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, actor)
	})

	h.router.PUT("/api/v1/actors/current/password", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)

		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		if !authenticatedActor.IsLocal {
			api.ErrorMessage(c, http.StatusForbidden, "not allowed")
			return
		}

		var request PasswordChangeRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.ChangePassword(ctx, authenticatedActor.Identifier, &request, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "password changed successfully")
	})

	h.router.PUT("/api/v1/actors/current/display-name", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		if !authenticatedActor.IsLocal {
			api.ErrorMessage(c, http.StatusForbidden, "not allowed")
			return
		}

		var request DisplayNameUpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.UpdateDisplayName(ctx, authenticatedActor.Identifier, &request, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "display name updated successfully")
	})

	h.router.PUT("/api/v1/actors/current/type", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		if !authenticatedActor.IsLocal {
			api.ErrorMessage(c, http.StatusForbidden, "not allowed")
			return
		}

		var request TypeUpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.UpdateType(ctx, authenticatedActor.Identifier, &request, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "actor type updated successfully")
	})

	h.router.DELETE("/api/v1/actors/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedActor, err := h.authenticator.ValidateContext(ctx, c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		if !authenticatedActor.IsLocal {
			api.ErrorMessage(c, http.StatusForbidden, "not allowed")
			return
		}

		err = h.manager.Delete(ctx, authenticatedActor.Identifier, authenticatedActor.TargetNodeIdentifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "actor deleted successfully")
	})
}
