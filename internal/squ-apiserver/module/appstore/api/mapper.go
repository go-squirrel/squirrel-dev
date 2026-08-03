package api

import (
	"squirrel-dev/internal/squ-apiserver/module/appstore/api/req"
	"squirrel-dev/internal/squ-apiserver/module/appstore/api/res"
	"squirrel-dev/internal/squ-apiserver/module/appstore/domain"
)

func toDomain(value req.AppStore) domain.App {
	return domain.App{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Category:    value.Category,
		Icon:        value.Icon,
		Version:     value.Version,
		Content:     value.Content,
		Tags:        value.Tags,
		Author:      value.Author,
		RepoURL:     value.RepoURL,
		HomepageURL: value.HomepageURL,
		IsOfficial:  value.IsOfficial,
		Downloads:   value.Downloads,
		Status:      value.Status,
	}
}

func fromDomain(value domain.App) res.AppStore {
	return res.AppStore{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Category:    value.Category,
		Icon:        value.Icon,
		Version:     value.Version,
		Content:     value.Content,
		Tags:        value.Tags,
		Author:      value.Author,
		RepoURL:     value.RepoURL,
		HomepageURL: value.HomepageURL,
		IsOfficial:  value.IsOfficial,
		Downloads:   value.Downloads,
		Status:      value.Status,
	}
}
