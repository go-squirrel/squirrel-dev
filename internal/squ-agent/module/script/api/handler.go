package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/script/api/req"
	"squirrel-dev/internal/squ-agent/module/script/application"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Execute(c *gin.Context) {
	value, ok := bindRequest[req.ExecuteScript](c)
	if !ok {
		return
	}
	err := h.service.Execute(c.Request.Context(), toExecuteRequest(value))
	writeResult(c, "Script execution task created", err)
}
