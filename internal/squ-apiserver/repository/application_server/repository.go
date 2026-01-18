package application_server

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(serverID uint) (applicationServers []model.ApplicationServer, err error)
	Get(id uint) (applicationServer model.ApplicationServer, err error)
	GetByServerAndApp(serverID, applicationID uint) (applicationServer model.ApplicationServer, err error)
	Delete(id uint) (err error)
	DeleteByApplicationID(applicationID uint) (err error)
	Add(req *model.ApplicationServer) (err error)
	Update(req *model.ApplicationServer) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
