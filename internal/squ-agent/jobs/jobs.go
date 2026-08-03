package jobs

import (
	"context"
	"strconv"
	"time"

	cronV3 "github.com/robfig/cron/v3"

	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-agent/config"
	applicationDomain "squirrel-dev/internal/squ-agent/module/application/domain"
	applicationInfra "squirrel-dev/internal/squ-agent/module/application/infra"
	configDomain "squirrel-dev/internal/squ-agent/module/config/domain"
	configInfra "squirrel-dev/internal/squ-agent/module/config/infra"
	monitorDomain "squirrel-dev/internal/squ-agent/module/monitor/domain"
	monitorInfra "squirrel-dev/internal/squ-agent/module/monitor/infra"
	scriptDomain "squirrel-dev/internal/squ-agent/module/script/domain"
	scriptInfra "squirrel-dev/internal/squ-agent/module/script/infra"
	"squirrel-dev/pkg/httpclient"
)

type HTTPPoster interface {
	Post(string, any, httpclient.Header) ([]byte, error)
}

type Jobs struct {
	config          *config.Config
	cron            *cronV3.Cron
	cache           cache.Cache
	applications    applicationDomain.Repository
	scriptTasks     scriptDomain.Repository
	configs         configDomain.Repository
	monitors        monitorDomain.Repository
	http            HTTPPoster
	composeProjects []composeProject
}

func New(
	conf *config.Config,
	cacheClient cache.Cache,
	agentDB, appDB, scriptTaskDB, monitorDB database.DB,
) *Jobs {
	return &Jobs{
		config:       conf,
		cron:         cronV3.New(cronV3.WithSeconds()),
		cache:        cacheClient,
		applications: applicationInfra.NewRepository(appDB.GetDB()),
		scriptTasks:  scriptInfra.NewRepository(scriptTaskDB.GetDB()),
		configs:      configInfra.NewRepository(agentDB.GetDB()),
		monitors:     monitorInfra.NewRepository(monitorDB.GetDB()),
		http:         httpclient.NewClient(10 * time.Second),
	}
}

func (j *Jobs) Start() error {
	if _, err := j.cron.AddFunc("*/5 * * * * *", j.checkApplicationStatus); err != nil {
		return err
	}
	if _, err := j.cron.AddFunc("*/5 * * * * *", j.reportScriptResults); err != nil {
		return err
	}
	if err := j.registerMonitorCollection(); err != nil {
		return err
	}
	if _, err := j.cron.AddFunc("*/5 * * * * *", j.refreshMonitorStatsCache); err != nil {
		return err
	}
	go j.refreshMonitorStatsCache()
	j.cron.Start()
	return nil
}

func (j *Jobs) configInt(ctx context.Context, key string) (int, error) {
	value, err := j.configs.GetByKey(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}
