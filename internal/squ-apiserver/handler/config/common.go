package config

import (
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"
	"squirrel-dev/internal/squ-apiserver/model"
)

func (c *Config) modelToResponse(daoC model.Config) res.Config {
	return res.Config{
		ID:    daoC.ID,
		Key:   daoC.Key,
		Value: daoC.Value,
	}
}

func (c *Config) requestToModel(request req.Config) model.Config {
	return model.Config{
		Key:   request.Key,
		Value: request.Value,
	}
}
