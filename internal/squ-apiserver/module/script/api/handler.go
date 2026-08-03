package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-apiserver/module/script/api/req"
	"squirrel-dev/internal/squ-apiserver/module/script/api/res"
	"squirrel-dev/internal/squ-apiserver/module/script/application"
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
	var result []res.Script
	for _, value := range values {
		result = append(result, toScriptResponse(value))
	}
	writeResult(c, result, err)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := scriptID(c)
	if !ok {
		return
	}
	value, err := h.service.Get(c.Request.Context(), id)
	writeResult(c, toScriptResponse(value), err)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := scriptID(c)
	if !ok {
		return
	}
	data, err := h.service.Delete(c.Request.Context(), id)
	writeResult(c, data, err)
}

func (h *Handler) Add(c *gin.Context) {
	request, ok := bindRequest[req.Script](c)
	if !ok {
		return
	}
	data, err := h.service.Add(c.Request.Context(), toScriptRequest(request))
	writeResult(c, data, err)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := scriptID(c)
	if !ok {
		return
	}
	request, ok := bindRequest[req.Script](c)
	if !ok {
		return
	}
	request.ID = id
	data, err := h.service.Update(c.Request.Context(), toScriptRequest(request))
	writeResult(c, data, err)
}

func (h *Handler) Execute(c *gin.Context) {
	request, ok := bindRequest[req.ExecuteScript](c)
	if !ok {
		return
	}
	data, err := h.service.Execute(c.Request.Context(), toExecuteRequest(request))
	writeResult(c, data, err)
}

func (h *Handler) ReceiveResult(c *gin.Context) {
	request, ok := bindRequest[req.ScriptResultReport](c)
	if !ok {
		return
	}
	data, err := h.service.ReceiveResult(c.Request.Context(), toResultReport(request))
	writeResult(c, data, err)
}

func (h *Handler) ListResults(c *gin.Context) {
	id, ok := scriptID(c)
	if !ok {
		return
	}
	values, err := h.service.ListResults(c.Request.Context(), id)
	var result []res.ScriptResult
	for _, value := range values {
		result = append(result, toScriptResultResponse(value))
	}
	writeResult(c, result, err)
}
