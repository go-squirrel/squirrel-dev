package cron

import (
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-agent/config"
	appRepository "squirrel-dev/internal/squ-agent/repository/application"
	confRepository "squirrel-dev/internal/squ-agent/repository/config"
	"squirrel-dev/internal/squ-agent/repository/monitor"
	scriptTaskRepo "squirrel-dev/internal/squ-agent/repository/script_task"
	"squirrel-dev/pkg/httpclient"
	"strconv"
	"time"

	cronV3 "github.com/robfig/cron/v3"
)

type Cron struct {
	Config         *config.Config
	Cron           *cronV3.Cron
	AppRepository  appRepository.Repository
	ScriptTaskRepo scriptTaskRepo.Repository
	ConfigRepo     confRepository.Repository
	MonitorRepo    monitor.Repository
	HTTPClient     *httpclient.Client
}

func New(config *config.Config, agentDB, appDB, scriptTaskDB, monitorDB database.DB) *Cron {

	c := cronV3.New(cronV3.WithSeconds())
	return &Cron{
		Config:         config,
		Cron:           c,
		AppRepository:  appRepository.New(appDB.GetDB()),
		ScriptTaskRepo: scriptTaskRepo.New(scriptTaskDB.GetDB()),
		ConfigRepo:     confRepository.New(agentDB.GetDB()),
		MonitorRepo:    monitor.New(monitorDB.GetDB()),
		HTTPClient:     httpclient.NewClient(10 * time.Second),
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

	err = c.startMonitor()
	if err != nil {
		return err
	}

	// 关键：启动 cron 定时器
	c.Cron.Start()
	return nil

}

func (c *Cron) getConfigByRepoToInt(key string) (int, error) {
	value, err := c.ConfigRepo.GetByKey(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (c *Cron) getConfigByRepo(key string) (string, error) {
	return c.ConfigRepo.GetByKey(key)
}
