package api

import (
	"squirrel-dev/internal/squ-apiserver/module/application/api/req"
	"squirrel-dev/internal/squ-apiserver/module/application/api/res"
	"squirrel-dev/internal/squ-apiserver/module/application/domain"
)

func toDomain(value req.Application) domain.Application {
	return domain.Application{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Content:     value.Content,
		Version:     value.Version,
	}
}

func fromDomain(value domain.Application) res.Application {
	return res.Application{
		ID:          value.ID,
		Name:        value.Name,
		Description: value.Description,
		Type:        value.Type,
		Content:     value.Content,
		Version:     value.Version,
	}
}
