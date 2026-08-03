package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("health"))
}
