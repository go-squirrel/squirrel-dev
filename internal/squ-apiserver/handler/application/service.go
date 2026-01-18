package application

import (
	"fmt"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/http"

	appRepository "squirrel-dev/internal/squ-apiserver/repository/application"
	appServerRepository "squirrel-dev/internal/squ-apiserver/repository/application_server"
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
)

type Application struct {
	Config            *config.Config
	Repository        appRepository.Repository
	AppServerRepo     appServerRepository.Repository
	ServerRepo        serverRepository.Repository
	HTTPClient        *http.Client
}

func New(config *config.Config, appRepo appRepository.Repository, appServerRepo appServerRepository.Repository, serverRepo serverRepository.Repository) *Application {
	return &Application{
		Config:        config,
		Repository:    appRepo,
		AppServerRepo: appServerRepo,
		ServerRepo:    serverRepo,
		HTTPClient:    http.NewClient(30 * time.Second),
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
			Status:      daoA.Status,
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
		Status:      daoA.Status,
		Content:     daoA.Content,
		Version:     daoA.Version,
	}

	return response.Success(appRes)
}

func (a *Application) Delete(id uint) response.Response {
	err := a.Repository.Delete(id)
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
		Status:      request.Status,
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
		Status:      request.Status,
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

// Deploy 部署应用到指定服务器
func (a *Application) Deploy(request req.DeployApplication) response.Response {
	// 1. 检查应用是否存在
	app, err := a.Repository.Get(request.ApplicationID)
	if err != nil {
		if err.Error() == "record not found" {
			return response.Error(res.ErrApplicationNotFound)
		}
		return response.Error(model.ReturnErrCode(err))
	}

	// 2. 检查服务器是否存在
	server, err := a.ServerRepo.Get(request.ServerID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 3. 检查是否已经部署到该服务器
	_, err = a.AppServerRepo.GetByServerAndApp(request.ServerID, request.ApplicationID)
	if err == nil {
		return response.Error(res.ErrAlreadyDeployed)
	}

	// 4. 发送部署请求到 agent
	agentURL := fmt.Sprintf("http://%s:%d/api/v1/application", server.IpAddress, server.AgentPort)
	agentReq := req.Application{
		Name:        app.Name,
		Description: app.Description,
		Type:        app.Type,
		Status:      app.Status,
		Content:     app.Content,
		Version:     app.Version,
	}

	_, err = a.HTTPClient.Post(agentURL, agentReq, nil)
	if err != nil {
		return response.Error(res.ErrDeployFailed)
	}

	// 5. 创建应用服务器关联记录
	appServer := model.ApplicationServer{
		ServerID:      request.ServerID,
		ApplicationID: request.ApplicationID,
	}

	err = a.AppServerRepo.Add(&appServer)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("deploy success")
}

