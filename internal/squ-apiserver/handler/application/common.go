package application

import (
	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
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

// returnApplicationErrCode 根据错误类型返回精确的应用错误码
func returnApplicationErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrApplicationNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrDuplicateApplication
	}
	return res.ErrApplicationUpdateFailed
}
