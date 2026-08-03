package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/monitor/application"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Stats(c *gin.Context) {
	value, err := h.service.Stats(c.Request.Context())
	writeResult(c, value, err)
}

func (h *Handler) DiskIO(c *gin.Context) {
	device, ok := monitorDevice(c)
	if !ok {
		return
	}
	value, err := h.service.DiskIO(c.Request.Context(), device)
	writeResult(c, value, err)
}

func (h *Handler) AllDiskIO(c *gin.Context) {
	value, err := h.service.AllDiskIO(c.Request.Context())
	writeResult(c, value, err)
}

func (h *Handler) NetIO(c *gin.Context) {
	interfaceName, ok := networkInterface(c)
	if !ok {
		return
	}
	value, err := h.service.NetIO(c.Request.Context(), interfaceName)
	writeResult(c, value, err)
}

func (h *Handler) AllNetIO(c *gin.Context) {
	value, err := h.service.AllNetIO(c.Request.Context())
	writeResult(c, value, err)
}

func (h *Handler) BaseByRange(c *gin.Context) {
	value, err := h.service.BaseByRange(c.Request.Context(), monitorTimeRange(c))
	writeResult(c, value, err)
}

func (h *Handler) DiskIOByRange(c *gin.Context) {
	value, err := h.service.DiskIOByRange(c.Request.Context(), monitorTimeRange(c))
	writeResult(c, value, err)
}

func (h *Handler) DiskUsageByRange(c *gin.Context) {
	value, err := h.service.DiskUsageByRange(c.Request.Context(), monitorTimeRange(c))
	writeResult(c, value, err)
}

func (h *Handler) NetworkByRange(c *gin.Context) {
	value, err := h.service.NetworkByRange(c.Request.Context(), monitorTimeRange(c))
	writeResult(c, value, err)
}
