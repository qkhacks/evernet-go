package node

import (
	"context"
	"github.com/evernetproto/evernet/internal/app/vertex/admin"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	router        *gin.Engine
	authenticator *admin.Authenticator
	manager       *Manager
}

func NewHandler(router *gin.Engine, authenticator *admin.Authenticator, manager *Manager) *Handler {
	return &Handler{router: router, authenticator: authenticator, manager: manager}
}

func (h *Handler) Register() {

	h.router.POST("/api/v1/nodes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedAdmin, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request CreationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		node, err := h.manager.Create(ctx, &request, authenticatedAdmin.Identifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, node)
	})

	h.router.GET("/api/v1/nodes", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		page, size := api.Page(c)

		nodes, err := h.manager.List(ctx, page, size)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, nodes)
	})

	h.router.GET("/api/v1/nodes/:identifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		identifier := c.Param("identifier")
		node, err := h.manager.Get(ctx, identifier)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, node)
	})

	h.router.PUT("/api/v1/nodes/:identifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request UpdateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		identifier := c.Param("identifier")
		err = h.manager.Update(ctx, identifier, &request)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "node updated successfully")
	})

	h.router.DELETE("/api/v1/nodes/:identifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("identifier")
		err = h.manager.Delete(ctx, identifier)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "node deleted successfully")
	})

	h.router.PUT("/api/v1/nodes/:identifier/signing-keys", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("identifier")
		response, err := h.manager.ResetSigningKeys(ctx, identifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
