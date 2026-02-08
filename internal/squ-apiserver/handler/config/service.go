package config

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"

	configRepository "squirrel-dev/internal/squ-apiserver/repository/config"

	"go.uber.org/zap"
)

type Config struct {
	Config     *config.Config
	Repository configRepository.Repository
}

func (c *Config) List() response.Response {
	var configs []res.Config
	daoConfigs, err := c.Repository.List()
	if err != nil {
		zap.L().Error("Failed to list configs", zap.Error(err))
		return response.Error(returnConfigErrCode(err))
	}
	for _, daoC := range daoConfigs {
		configs = append(configs, c.modelToResponse(daoC))
	}
	return response.Success(configs)
}

func (c *Config) Get(id uint) response.Response {
	daoC, err := c.Repository.Get(id)
	if err != nil {
		zap.L().Error("Failed to get config", zap.Uint("id", id), zap.Error(err))
		return response.Error(returnConfigErrCode(err))
	}
	configRes := c.modelToResponse(daoC)

	return response.Success(configRes)
}

func (c *Config) Delete(id uint) response.Response {
	err := c.Repository.Delete(id)
	if err != nil {
		zap.L().Error("Failed to delete config", zap.Uint("id", id), zap.Error(err))
		return response.Error(returnConfigErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Add(request req.Config) response.Response {
	modelReq := c.requestToModel(request)

	err := c.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("Failed to add config", zap.String("key", request.Key), zap.Error(err))
		return response.Error(returnConfigErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Update(request req.Config) response.Response {
	modelReq := c.requestToModel(request)
	modelReq.ID = request.ID
	err := c.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("Failed to update config", zap.Uint("id", request.ID), zap.String("key", request.Key), zap.Error(err))
		return response.Error(returnConfigErrCode(err))
	}

	return response.Success("success")
}
