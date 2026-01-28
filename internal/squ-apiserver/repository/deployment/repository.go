package deployment

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(serverID uint) (applicationServers []model.Deployment, err error)
	Get(id uint) (applicationServer model.Deployment, err error)
	GetByServerAndApp(serverID, applicationID uint) (applicationServer model.Deployment, err error)
	Delete(id uint) (err error)
	DeleteByApplicationID(applicationID uint) (err error)
	Add(req *model.Deployment) (err error)
	Update(req *model.Deployment) (err error)
	UpdateStatus(serverID, applicationID uint, status string) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
