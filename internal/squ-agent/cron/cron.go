package cron

import (
	"squirrel-dev/internal/pkg/database"
	appRepository "squirrel-dev/internal/squ-agent/repository/application"
	"squirrel-dev/internal/squ-agent/repository/config"
	"squirrel-dev/internal/squ-agent/repository/monitor"
	scriptTaskRepo "squirrel-dev/internal/squ-agent/repository/script_task"

	cronV3 "github.com/robfig/cron/v3"
)

type Cron struct {
	Cron           *cronV3.Cron
	AppRepository  appRepository.Repository
	ScriptTaskRepo scriptTaskRepo.Repository
	ConfigRepo     config.Repository
	MonitorRepo    monitor.Repository
}

func New(agentDB, appDB, scriptTaskDB, monitorDB database.DB) *Cron {

	c := cronV3.New(cronV3.WithSeconds())
	return &Cron{
		Cron:           c,
		AppRepository:  appRepository.New(appDB.GetDB()),
		ScriptTaskRepo: scriptTaskRepo.New(scriptTaskDB.GetDB()),
		ConfigRepo:     config.New(agentDB.GetDB()),
		MonitorRepo:    monitor.New(monitorDB.GetDB()),
	}
}

func (c *Cron) Start() error {
	err := c.startApp()
	if err != nil {
		return err
	}

	err = c.startScriptResultReporter()
	if err != nil {
		return err
	}

	err = c.startMonitor(c.ConfigRepo, c.MonitorRepo)
	if err != nil {
		return err
	}

	// 关键：启动 cron 定时器
	c.Cron.Start()
	return nil

}
