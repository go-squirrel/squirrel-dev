package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-apiserver/module/deployment/api/req"
	"squirrel-dev/internal/squ-apiserver/module/deployment/api/res"
	"squirrel-dev/internal/squ-apiserver/module/deployment/application"
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
	serverID, ok := deploymentServerID(c)
	if !ok {
		return
	}
	values, err := h.service.List(c.Request.Context(), serverID)
	var result []res.Deployment
	for _, value := range values {
		result = append(result, toDeploymentResponse(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := deploymentID(c)
	if !ok {
		return
	}
	request, ok := bindRequest[req.UpdateDeployment](c)
	if !ok {
		return
	}
	data, err := h.service.Update(c.Request.Context(), id, request.Content)
	writeResult(c, data, err)
}

func (h *Handler) Deploy(c *gin.Context) {
	id, ok := applicationID(c)
	if !ok {
		return
	}
	request, ok := bindRequest[req.DeployApplication](c)
	if !ok {
		return
	}
	data, err := h.service.Deploy(c.Request.Context(), id, request.ServerID)
	writeResult(c, data, err)
}

func (h *Handler) ListServers(c *gin.Context) {
	id, ok := applicationID(c)
	if !ok {
		return
	}
	values, err := h.service.ListServers(c.Request.Context(), id)
	var result []res.ServerInfo
	for _, value := range values {
		result = append(result, toServerResponse(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Undeploy(c *gin.Context) {
	id, ok := deploymentID(c)
	if !ok {
		return
	}
	data, err := h.service.Undeploy(c.Request.Context(), id)
	writeResult(c, data, err)
}

func (h *Handler) Stop(c *gin.Context) {
	id, ok := deploymentID(c)
	if !ok {
		return
	}
	data, err := h.service.Stop(c.Request.Context(), id)
	writeResult(c, data, err)
}

func (h *Handler) Start(c *gin.Context) {
	id, ok := deploymentID(c)
	if !ok {
		return
	}
	data, err := h.service.Start(c.Request.Context(), id)
	writeResult(c, data, err)
}

func (h *Handler) ReDeploy(c *gin.Context) {
	id, ok := deploymentID(c)
	if !ok {
		return
	}
	data, err := h.service.ReDeploy(c.Request.Context(), id)
	writeResult(c, data, err)
}

func (h *Handler) ReportStatus(c *gin.Context) {
	request, ok := bindRequest[req.ReportApplicationStatus](c)
	if !ok {
		return
	}
	data, err := h.service.ReportStatus(c.Request.Context(), request.DeployID, request.Status)
	writeResult(c, data, err)
}
