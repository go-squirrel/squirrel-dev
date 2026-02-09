package deployment

import (
	"squirrel-dev/internal/squ-apiserver/handler/deployment/res"

	"gorm.io/gorm"
)

// returnDeploymentErrCode 根据错误类型返回精确的部署错误码
func returnDeploymentErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrDeploymentNotFound
	case gorm.ErrDuplicatedKey:
		return res.ErrAlreadyDeployed
	}
	return res.ErrCreateDeploymentRecordFailed
}
