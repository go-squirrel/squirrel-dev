package application

import (
	"fmt"

	"squirrel-dev/internal/squ-apiserver/handler/application/req"
	"squirrel-dev/internal/squ-apiserver/handler/application/res"
	"squirrel-dev/internal/squ-apiserver/model"

	"gopkg.in/yaml.v3"
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

// validateYAML 验证 YAML 格式是否正确
func validateYAML(content string) error {
	if content == "" {
		return nil // 空内容视为有效
	}
	var result any
	if err := yaml.Unmarshal([]byte(content), &result); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}
	return nil
}
