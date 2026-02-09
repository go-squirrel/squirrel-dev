package config

import (
	"squirrel-dev/internal/squ-apiserver/handler/config/req"
	"squirrel-dev/internal/squ-apiserver/handler/config/res"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
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

// returnConfigErrCode 根据错误类型返回精确的配置错误码
func returnConfigErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrConfigNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrConfigKeyAlreadyExists
	}
	return res.ErrConfigUpdateFailed
}
