package application

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"

	"squirrel-dev/internal/squ-apiserver/handler/auth/req"
	authModel "squirrel-dev/internal/squ-apiserver/model/auth"
)

type Application struct {
	Config      *config.Config
	ModelClient authModel.Repository
}

func (a *Application) List(request req.Request) response.Response {

	return response.Success("")
}
