package application

import (
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"
)

func (a *Application) modelToResponse(daoA model.Application) res.Application {
	return res.Application{
		ID:          daoA.ID,
		Name:        daoA.Name,
		Description: daoA.Description,
		Type:        daoA.Type,
		Content:     daoA.Content,
		Version:     daoA.Version,
	}
}

func (a *Application) requestToModel(request req.Application) model.Application {
	return model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Content:     request.Content,
		Version:     request.Version,
	}
}
