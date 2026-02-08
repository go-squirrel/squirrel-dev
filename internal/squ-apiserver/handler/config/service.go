package config

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"
	"squirrel-dev/internal/squ-apiserver/model"

	configRepository "squirrel-dev/internal/squ-apiserver/repository/config"
)

type Config struct {
	Config     *config.Config
	Repository configRepository.Repository
}

func (c *Config) List() response.Response {
	var configs []res.Config
	daoConfigs, err := c.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoC := range daoConfigs {
		configs = append(configs, c.modelToResponse(daoC))
	}
	return response.Success(configs)
}

func (c *Config) Get(id uint) response.Response {
	daoC, err := c.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	configRes := c.modelToResponse(daoC)

	return response.Success(configRes)
}

func (c *Config) Delete(id uint) response.Response {
	err := c.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Add(request req.Config) response.Response {
	modelReq := c.requestToModel(request)

	err := c.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (c *Config) Update(request req.Config) response.Response {
	modelReq := c.requestToModel(request)
	modelReq.ID = request.ID
	err := c.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
