package deployment

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(serverID uint) (dms []model.Deployment, err error)
	Get(id uint) (dm model.Deployment, err error)
	GetByDeployID(deployID uint64) (dm model.Deployment, err error)
	GetByServerAndApp(serverID, applicationID uint) (dm model.Deployment, err error)
	Delete(id uint) (err error)
	DeleteByApplicationID(applicationID uint) (err error)
	Add(req *model.Deployment) (err error)
	Update(req *model.Deployment) (err error)
	UpdateStatus(deployID uint64, status string) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
