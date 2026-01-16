package app_store

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/req"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/res"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/compose"

	appStoreRepository "squirrel-dev/internal/squ-apiserver/repository/app_store"
)

type AppStore struct {
	Config     *config.Config
	Repository appStoreRepository.Repository
}

func (a *AppStore) List() response.Response {
	var appStores []res.AppStore
	daoAppStores, err := a.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoA := range daoAppStores {
		appStores = append(appStores, res.AppStore{
			ID:          daoA.ID,
			Name:        daoA.Name,
			Description: daoA.Description,
			Type:        daoA.Type,
			Category:    daoA.Category,
			Icon:        daoA.Icon,
			Version:     daoA.Version,
			Content:     daoA.Content,
			Tags:        daoA.Tags,
			Author:      daoA.Author,
			RepoUrl:     daoA.RepoUrl,
			HomepageUrl: daoA.HomepageUrl,
			IsOfficial:  daoA.IsOfficial,
			Downloads:   daoA.Downloads,
			Status:      daoA.Status,
		})
	}
	return response.Success(appStores)
}

func (a *AppStore) Get(id uint) response.Response {
	var appStoreRes res.AppStore
	daoA, err := a.Repository.Get(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}
	appStoreRes = res.AppStore{
		ID:          daoA.ID,
		Name:        daoA.Name,
		Description: daoA.Description,
		Type:        daoA.Type,
		Category:    daoA.Category,
		Icon:        daoA.Icon,
		Version:     daoA.Version,
		Content:     daoA.Content,
		Tags:        daoA.Tags,
		Author:      daoA.Author,
		RepoUrl:     daoA.RepoUrl,
		HomepageUrl: daoA.HomepageUrl,
		IsOfficial:  daoA.IsOfficial,
		Downloads:   daoA.Downloads,
		Status:      daoA.Status,
	}

	return response.Success(appStoreRes)
}

func (a *AppStore) Delete(id uint) response.Response {
	err := a.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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
		if err := compose.ValidateContent(request.Content); err != nil {
			return response.ErrorUnknown(res.ErrInvalidComposeContent, err.Error())
		}
	}

	modelReq := model.AppStore{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Category:    request.Category,
		Icon:        request.Icon,
		Version:     request.Version,
		Content:     request.Content,
		Tags:        request.Tags,
		Author:      request.Author,
		RepoUrl:     request.RepoUrl,
		HomepageUrl: request.HomepageUrl,
		IsOfficial:  request.IsOfficial,
		Downloads:   request.Downloads,
		Status:      request.Status,
	}

	err := a.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
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
		if err := compose.ValidateContent(request.Content); err != nil {
			return response.ErrorUnknown(res.ErrInvalidComposeContent, err.Error())
		}
	}

	modelReq := model.AppStore{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Category:    request.Category,
		Icon:        request.Icon,
		Version:     request.Version,
		Content:     request.Content,
		Tags:        request.Tags,
		Author:      request.Author,
		RepoUrl:     request.RepoUrl,
		HomepageUrl: request.HomepageUrl,
		IsOfficial:  request.IsOfficial,
		Downloads:   request.Downloads,
		Status:      request.Status,
	}
	modelReq.ID = request.ID
	err := a.Repository.Update(&modelReq)

	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
