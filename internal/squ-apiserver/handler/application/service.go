package application

import (
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"

	appRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	appServerRepository "squirrel-dev/internal/squ-apiserver/repository/application_server"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"

	"go.uber.org/zap"
)

type Application struct {
	Config        *config.Config
	Repository    appRepository.Repository
	AppServerRepo appServerRepository.Repository
	ServerRepo    serverRepository.Repository
	HTTPClient    *httpclient.Client
}

func New(config *config.Config, appRepo appRepository.Repository, appServerRepo appServerRepository.Repository, serverRepo serverRepository.Repository) *Application {
	hc := httpclient.NewClient(30 * time.Second)
	return &Application{
		Config:        config,
		Repository:    appRepo,
		AppServerRepo: appServerRepo,
		ServerRepo:    serverRepo,
		HTTPClient:    hc,
	}
}

func (a *Application) List() response.Response {
	var applications []res.Application
	daoApps, err := a.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoA := range daoApps {
		applications = append(applications, res.Application{
			ID:          daoA.ID,
			Name:        daoA.Name,
			Description: daoA.Description,
			Type:        daoA.Type,
			Content:     daoA.Content,
			Version:     daoA.Version,
		})
	}
	return response.Success(applications)
}

func (a *Application) Get(id uint) response.Response {
	var appRes res.Application
	daoA, err := a.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	appRes = res.Application{
		ID:          daoA.ID,
		Name:        daoA.Name,
		Description: daoA.Description,
		Type:        daoA.Type,
		Content:     daoA.Content,
		Version:     daoA.Version,
	}

	return response.Success(appRes)
}

func (a *Application) Delete(id uint) response.Response {
	// 先删除应用服务器关联记录
	err := a.AppServerRepo.DeleteByApplicationID(id)
	if err != nil {
		zap.L().Error("删除应用服务器关联记录失败",
			zap.Uint("application_id", id),
			zap.Error(err),
		)
		// 不返回错误，继续删除应用记录
	}

	// 删除应用记录
	err = a.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Add(request req.Application) response.Response {
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Content:     request.Content,
		Version:     request.Version,
	}

	err := a.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Update(request req.Application) response.Response {
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Content:     request.Content,
		Version:     request.Version,
	}
	modelReq.ID = request.ID
	err := a.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) ReportStatus(request req.ReportApplicationStatus) response.Response {
	// 验证应用服务器关联记录是否存在
	_, err := a.AppServerRepo.GetByServerAndApp(request.ServerID, request.ApplicationID)
	if err != nil {
		zap.L().Error("应用服务器关联记录不存在",
			zap.Uint("server_id", request.ServerID),
			zap.Uint("application_id", request.ApplicationID),
			zap.Error(err),
		)
		return response.Error(response.ErrCodeParameter)
	}

	// 更新状态
	err = a.AppServerRepo.UpdateStatus(request.ServerID, request.ApplicationID, request.Status)
	if err != nil {
		zap.L().Error("更新应用状态失败",
			zap.Uint("server_id", request.ServerID),
			zap.Uint("application_id", request.ApplicationID),
			zap.String("status", request.Status),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("应用状态已更新",
		zap.Uint("server_id", request.ServerID),
		zap.Uint("application_id", request.ApplicationID),
		zap.String("status", request.Status),
	)

	return response.Success("success")
}
