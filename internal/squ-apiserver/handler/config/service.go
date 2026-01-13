package config

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"

	"squirrel-dev/internal/squ-apiserver/handler/auth/req"
	authModel "squirrel-dev/internal/squ-apiserver/model/auth"
)

type Config struct {
	Config      *config.Config
	ModelClient authModel.Repository
}

func (c *Config) List(request req.Request) response.Response {

	return response.Success("")
}
