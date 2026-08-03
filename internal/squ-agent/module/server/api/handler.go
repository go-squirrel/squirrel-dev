package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/server/application"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Info(c *gin.Context) {
	value, err := h.service.GetInfo(c.Request.Context())
	writeResult(c, toResponse(value), err)
}
