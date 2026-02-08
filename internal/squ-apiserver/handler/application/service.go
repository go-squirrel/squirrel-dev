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
	serverRepository "squirrel-dev/internal/squ-apiserver/repository/server"
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
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoA := range daoApps {
		applications = append(applications, a.modelToResponse(daoA))
	}
	return response.Success(applications)
}

func (a *Application) Get(id uint) response.Response {
	daoA, err := a.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	appRes := a.modelToResponse(daoA)

	return response.Success(appRes)
}

func (a *Application) Delete(id uint) response.Response {

	// 删除应用记录
	err := a.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Add(request req.Application) response.Response {
	modelReq := a.requestToModel(request)

	err := a.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Update(request req.Application) response.Response {
	modelReq := a.requestToModel(request)
	modelReq.ID = request.ID
	err := a.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
