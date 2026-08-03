package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/squ-apiserver/module/monitor/application"
	"squirrel-dev/internal/squ-apiserver/module/monitor/domain"
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
	h.writeForServer(c, h.service.Stats)
}

func (h *Handler) DiskIO(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		device := c.Param("device")
		if device == "" {
			zap.L().Warn("monitor device parameter is empty", zap.Uint("server_id", id))
			return domain.Result{}, errInvalidConfig
		}
		return h.service.DiskIO(id, device)
	})
}

func (h *Handler) AllDiskIO(c *gin.Context) {
	h.writeForServer(c, h.service.AllDiskIO)
}

func (h *Handler) NetIO(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		interfaceName := c.Param("interface")
		if interfaceName == "" {
			zap.L().Warn("monitor interface parameter is empty", zap.Uint("server_id", id))
			return domain.Result{}, errInvalidConfig
		}
		return h.service.NetIO(id, interfaceName)
	})
}

func (h *Handler) AllNetIO(c *gin.Context) {
	h.writeForServer(c, h.service.AllNetIO)
}

func (h *Handler) BaseRange(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		return h.service.BaseRange(id, c.DefaultQuery("range", "1h"))
	})
}

func (h *Handler) DiskRange(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		return h.service.DiskRange(id, c.DefaultQuery("range", "1h"))
	})
}

func (h *Handler) DiskUsageRange(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		return h.service.DiskUsageRange(id, c.DefaultQuery("range", "1h"))
	})
}

func (h *Handler) NetworkRange(c *gin.Context) {
	h.writeForServer(c, func(id uint) (domain.Result, error) {
		return h.service.NetworkRange(id, c.DefaultQuery("range", "1h"))
	})
}
