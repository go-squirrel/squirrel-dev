package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-apiserver/module/server/api/req"
	"squirrel-dev/internal/squ-apiserver/module/server/api/res"
	"squirrel-dev/internal/squ-apiserver/module/server/application"
)

type Handler struct {
	service    *application.Service
	signingKey string
}

func NewHandler(service *application.Service, signingKey string) *Handler {
	return &Handler{
		service:    service,
		signingKey: signingKey,
	}
}

func (h *Handler) List(c *gin.Context) {
	values, err := h.service.List(c.Request.Context())
	var result []res.Server
	for _, value := range values {
		result = append(result, toResponse(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := serverID(c)
	if !ok {
		return
	}
	value, err := h.service.Get(c.Request.Context(), id)
	writeResult(c, toResponse(value), err)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := serverID(c)
	if !ok {
		return
	}
	err := h.service.Delete(c.Request.Context(), id)
	writeResult(c, "success", err)
}

func (h *Handler) Add(c *gin.Context) {
	request, ok := bindRequest[req.Server](c)
	if !ok {
		return
	}
	err := h.service.Add(c.Request.Context(), toApplication(request))
	writeResult(c, "success", err)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := serverID(c)
	if !ok {
		return
	}
	request, ok := bindRequest[req.Server](c)
	if !ok {
		return
	}
	request.ID = id
	err := h.service.Update(c.Request.Context(), toApplication(request))
	writeResult(c, "success", err)
}

func (h *Handler) TestSSH(c *gin.Context) {
	id, ok := serverID(c)
	if !ok {
		return
	}
	server, err := h.service.TestSSH(c.Request.Context(), id)
	writeSSHResult(c, toSSHTestResponse(server), err)
}

func (h *Handler) CheckAgent(c *gin.Context) {
	request, ok := bindRequest[req.CheckAgent](c)
	if !ok {
		return
	}
	ready, message, info := h.service.CheckAgent(c.Request.Context(), request.IPAddress, request.Port)
	writeResult(c, toAgentCheckResponse(ready, message, info), nil)
}
