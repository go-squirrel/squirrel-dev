package config

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/config/req"
	"squirrel-dev/internal/squ-agent/handler/config/res"
	"squirrel-dev/internal/squ-agent/model"

	configRepository "squirrel-dev/internal/squ-agent/repository/config"

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
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoC := range daoConfigs {
		configs = append(configs, res.Config{
			ID:    daoC.ID,
			Key:   daoC.Key,
			Value: daoC.Value,
		})
	}
	return response.Success(configs)
}

func (c *Config) Get(id uint) response.Response {
	var configRes res.Config
	daoC, err := c.Repository.Get(id)
	if err != nil {
		zap.L().Error("Failed to get config", zap.Uint("id", id), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}
	configRes = res.Config{
		ID:    daoC.ID,
		Key:   daoC.Key,
		Value: daoC.Value,
	}

	return response.Success(configRes)
}

func (c *Config) Delete(id uint) response.Response {
	err := c.Repository.Delete(id)
	if err != nil {
		zap.L().Error("Failed to delete config", zap.Uint("id", id), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Add(request req.Config) response.Response {
	modelReq := model.Config{
		Key:   request.Key,
		Value: request.Value,
	}

	err := c.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("Failed to add config", zap.String("key", request.Key), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Update(request req.Config) response.Response {
	modelReq := model.Config{
		Key:   request.Key,
		Value: request.Value,
	}
	modelReq.ID = request.ID
	err := c.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("Failed to update config", zap.Uint("id", request.ID), zap.String("key", request.Key), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
