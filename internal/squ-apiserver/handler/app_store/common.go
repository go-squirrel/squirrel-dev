package app_store

import (
	"squirrel-dev/internal/squ-apiserver/handler/app_store/req"
	"squirrel-dev/internal/squ-apiserver/handler/app_store/res"
	"squirrel-dev/internal/squ-apiserver/model"
)

func (a *AppStore) modelToResponse(daoA model.AppStore) res.AppStore {
	return res.AppStore{
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
}

func (a *AppStore) requestToModel(request req.AppStore) model.AppStore {
	return model.AppStore{
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
}
