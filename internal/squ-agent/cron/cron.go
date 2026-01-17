package cron

import (
	"squirrel-dev/internal/pkg/database"
	appRepository "squirrel-dev/internal/squ-agent/repository/application"

	cronV3 "github.com/robfig/cron/v3"
)

type Cron struct {
	Cron          *cronV3.Cron
	AppRepository appRepository.Repository
}

func New(appDB database.DB) *Cron {
	c := cronV3.New(cronV3.WithParser(cronV3.NewParser(
		cronV3.SecondOptional | cronV3.Minute | cronV3.Hour | cronV3.Dom | cronV3.Month | cronV3.Dow | cronV3.Descriptor,
	)))
	return &Cron{
		Cron:          c,
		AppRepository: appRepository.New(appDB.GetDB()),
	}
}

func (c *Cron) Start() {
	c.startApp()
}
