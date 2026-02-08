package app_store

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/req"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/res"
	"squirrel-dev/pkg/compose"

	appStoreRepository "squirrel-dev/internal/squ-apiserver/repository/app_store"

	"go.uber.org/zap"
)

type AppStore struct {
	Config     *config.Config
	Repository appStoreRepository.Repository
}

func (a *AppStore) List() response.Response {
	var appStores []res.AppStore
	daoAppStores, err := a.Repository.List()
	if err != nil {
		return response.Error(returnAppStoreErrCode(err))
	}
	for _, daoA := range daoAppStores {
		appStores = append(appStores, a.modelToResponse(daoA))
	}
	return response.Success(appStores)
}

func (a *AppStore) Get(id uint) response.Response {
	var appStoreRes res.AppStore
	daoA, err := a.Repository.Get(id)
	if err != nil {
		return response.Error(returnAppStoreErrCode(err))
	}
	appStoreRes = a.modelToResponse(daoA)

	return response.Success(appStoreRes)
}

func (a *AppStore) Delete(id uint) response.Response {
	err := a.Repository.Delete(id)
	if err != nil {
		return response.Error(returnAppStoreErrCode(err))
	}

	return response.Success("success")
}

func (a *AppStore) Add(request req.AppStore) response.Response {
	// 验证应用类型
	if !compose.IsValidAppType(request.Type) {
		return response.Error(res.ErrUnsupportedAppType)
	}

	// 如果是 compose 类型，验证 compose 内容
	if request.Type == "compose" {
		request.Content = compose.TrimSpaceContent(request.Content)
		if err := compose.ValidateContent(request.Name, request.Content); err != nil {
			zap.S().Error(err)
			return response.Error(res.ErrInvalidComposeContent)
		}
	}

	modelReq := a.requestToModel(request)

	err := a.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(returnAppStoreErrCode(err))
	}

	return response.Success("success")
}

func (a *AppStore) Update(request req.AppStore) response.Response {
	// 验证应用类型
	if !compose.IsValidAppType(request.Type) {
		return response.Error(res.ErrUnsupportedAppType)
	}

	// 如果是 compose 类型，验证 compose 内容
	if request.Type == "compose" {
		request.Content = compose.TrimSpaceContent(request.Content)
		if err := compose.ValidateContent(request.Name, request.Content); err != nil {
			zap.S().Error(err)
			return response.Error(res.ErrInvalidComposeContent)
		}
	}

	modelReq := a.requestToModel(request)
	modelReq.ID = request.ID
	err := a.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(returnAppStoreErrCode(err))
	}

	return response.Success("success")
}
