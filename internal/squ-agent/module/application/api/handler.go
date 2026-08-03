package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-agent/module/application/api/res"
	appService "squirrel-dev/internal/squ-agent/module/application/application"
)

type Handler struct {
	service *appService.Service
}

func NewHandler(service *appService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c *gin.Context) {
	values, err := h.service.List(c.Request.Context())
	var result []res.Application
	for _, value := range values {
		result = append(result, fromDomain(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := applicationID(c)
	if !ok {
		return
	}
	value, err := h.service.Get(c.Request.Context(), id)
	writeResult(c, fromDomain(value), err)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := applicationID(c)
	if !ok {
		return
	}
	writeApplicationResult(c, h.service.Delete(c.Request.Context(), id))
}

func (h *Handler) Add(c *gin.Context) {
	value, ok := bindApplication(c)
	if !ok {
		return
	}
	writeApplicationResult(c, h.service.Add(c.Request.Context(), toApplicationRequest(value)))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := applicationID(c)
	if !ok {
		return
	}
	value, ok := bindApplication(c)
	if !ok {
		return
	}
	value.ID = id
	err := h.service.Update(c.Request.Context(), toApplicationRequest(value))
	writeResult(c, "success", err)
}

func (h *Handler) Start(c *gin.Context) {
	deployID, ok := deploymentID(c)
	if !ok {
		return
	}
	writeApplicationResult(c, h.service.Start(c.Request.Context(), deployID))
}

func (h *Handler) Stop(c *gin.Context) {
	deployID, ok := deploymentID(c)
	if !ok {
		return
	}
	writeApplicationResult(c, h.service.Stop(c.Request.Context(), deployID))
}

func (h *Handler) DeleteByDeployID(c *gin.Context) {
	deployID, ok := deploymentID(c)
	if !ok {
		return
	}
	writeApplicationResult(c, h.service.DeleteByDeployID(c.Request.Context(), deployID))
}
