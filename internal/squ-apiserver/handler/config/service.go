package config

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"
	"squirrel-dev/internal/squ-apiserver/model"

	configModel "squirrel-dev/internal/squ-apiserver/model/config"
)

type Config struct {
	Config      *config.Config
	ModelClient configModel.Repository
}

func (c *Config) List() response.Response {
	var configs []res.Config
	daoConfigs, err := c.ModelClient.List()
	if err != nil {
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
	daoC, err := c.ModelClient.Get(id)
	if err != nil {
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
	err := c.ModelClient.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Add(request req.Config) response.Response {
	modelReq := configModel.Config{
		Key:   request.Key,
		Value: request.Value,
	}

	err := c.ModelClient.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Update(request req.Config) response.Response {
	modelReq := configModel.Config{
		Key:   request.Key,
		Value: request.Value,
	}
	modelReq.ID = request.ID
	err := c.ModelClient.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

