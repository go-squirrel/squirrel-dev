package application

import (
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/pkg/httpclient"

	appRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

type Application struct {
	Config     *config.Config
	Repository appRepository.Repository
	ServerRepo serverRepository.Repository
	HTTPClient *httpclient.Client
}

func New(config *config.Config, appRepo appRepository.Repository, serverRepo serverRepository.Repository) *Application {
	hc := httpclient.NewClient(30 * time.Second)
	return &Application{
		Config:     config,
		Repository: appRepo,
		ServerRepo: serverRepo,
		HTTPClient: hc,
	}
}

func (a *Application) List() response.Response {
	var applications []res.Application
	daoApps, err := a.Repository.List()
	if err != nil {
		zap.L().Error("Failed to list applications", zap.Error(err))
		return response.Error(returnApplicationErrCode(err))
	}
	for _, daoA := range daoApps {
		applications = append(applications, a.modelToResponse(daoA))
	}
	return response.Success(applications)
}

func (a *Application) Get(id uint) response.Response {
	daoA, err := a.Repository.Get(id)
	if err != nil {
		zap.L().Error("Failed to get application", zap.Uint("id", id), zap.Error(err))
		return response.Error(returnApplicationErrCode(err))
	}
	appRes := a.modelToResponse(daoA)

	return response.Success(appRes)
}

func (a *Application) Delete(id uint) response.Response {

	// 删除应用记录
	err := a.Repository.Delete(id)
	if err != nil {
		zap.L().Error("Failed to delete application", zap.Uint("id", id), zap.Error(err))
		return response.Error(returnApplicationErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Add(request req.Application) response.Response {
	modelReq := a.requestToModel(request)

	// 验证 Content 的 YAML 格式
	if err := validateYAML(modelReq.Content); err != nil {
		zap.L().Error("Invalid YAML content", zap.String("name", request.Name), zap.Error(err))
		return response.Error(res.ErrInvalidApplicationConfig)
	}

	err := a.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("Failed to add application", zap.String("name", request.Name), zap.Error(err))
		return response.Error(returnApplicationErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Update(request req.Application) response.Response {
	modelReq := a.requestToModel(request)
	modelReq.ID = request.ID

	// 验证 Content 的 YAML 格式
	if err := validateYAML(modelReq.Content); err != nil {
		zap.L().Error("Invalid YAML content", zap.Uint("id", request.ID), zap.String("name", request.Name), zap.Error(err))
		return response.Error(res.ErrInvalidApplicationConfig)
	}

	err := a.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("Failed to update application", zap.Uint("id", request.ID), zap.String("name", request.Name), zap.Error(err))
		return response.Error(returnApplicationErrCode(err))
	}

	return response.Success("success")
}
