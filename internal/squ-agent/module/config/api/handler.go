package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/config/api/res"
	"squirrel-dev/internal/squ-agent/module/config/application"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c *gin.Context) {
	values, err := h.service.List(c.Request.Context())
	var result []res.Config
	for _, value := range values {
		result = append(result, fromDomain(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := configID(c)
	if !ok {
		return
	}
	value, err := h.service.Get(c.Request.Context(), id)
	writeResult(c, fromDomain(value), err)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := configID(c)
	if !ok {
		return
	}
	err := h.service.Delete(c.Request.Context(), id)
	writeResult(c, "success", err)
}

func (h *Handler) Save(c *gin.Context) {
	value, ok := bindConfig(c)
	if !ok {
		return
	}
	err := h.service.Save(c.Request.Context(), value.Key, value.Value)
	writeResult(c, "success", err)
}
