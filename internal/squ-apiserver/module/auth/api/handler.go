package api

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/squ-apiserver/module/auth/application"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	request, ok := bindLogin(c)
	if !ok {
		return
	}
	token, err := h.service.Login(c.Request.Context(), request.Username, request.Password)
	writeResult(c, toTokenResponse(token), err)
}
