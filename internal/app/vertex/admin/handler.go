package admin

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

	h.router.GET("/api/v1/admins/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedAdmin, err := h.authenticator.ValidateContext(c)

		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		admin, err := h.manager.Get(ctx, authenticatedAdmin.Identifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, admin)
	})

	h.router.PUT("/api/v1/admins/current/password", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedAdmin, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request PasswordChangeRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		err = h.manager.ChangePassword(ctx, authenticatedAdmin.Identifier, &request)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		api.Success(c, http.StatusOK, "admin password changed successfully")
	})

	h.router.POST("/api/v1/admins", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		authenticatedAdmin, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		var request AdditionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			api.Error(c, http.StatusBadRequest, err)
			return
		}

		admin, err := h.manager.Add(ctx, &request, authenticatedAdmin.Identifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, admin)
	})

	h.router.GET("/api/v1/admins/:identifier", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		identifier := c.Param("identifier")

		admin, err := h.manager.Get(ctx, identifier)

		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, admin)
	})

	h.router.DELETE("/api/v1/admins/:identifier", func(c *gin.Context) {
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

		api.Success(c, http.StatusOK, "admin deleted successfully")
	})

	h.router.GET("/api/v1/admins", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateContext(c)
		if err != nil {
			api.Error(c, http.StatusUnauthorized, err)
			return
		}

		page, size := api.Page(c)

		admins, err := h.manager.List(ctx, page, size)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, admins)
	})
}
