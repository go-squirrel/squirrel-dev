package options

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/log"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/cron"
	"squirrel-dev/internal/squ-agent/server"
)

type AppOptions struct {
	ConfFile string
	Config   *config.Config
}

func NewAppOptions() *AppOptions {
	o := &AppOptions{}
	return o
}

func (o *AppOptions) NewServer() (*server.Server, error) {
	s := server.NewServer()
	o.loadConfig(o.ConfFile)
	s.Config = o.Config

	gin.SetMode(s.Config.Server.Mode)
	s.Gin = gin.New()

	s.Log = log.NewClient(o.Config.Log.InfoFilePath, o.Config.Log.ErrorFilePath, o.Config.Log.Level,
		o.Config.Log.MaxSize, o.Config.Log.MaxBackups, o.Config.Log.MaxAge,
	)

	if o.Config.DB.Type == "sqlite" {
		s.AgentDB = database.New(o.Config.DB.Type, o.Config.DB.Sqlite.AgentFilePath, database.WithMigrate(true))
		s.AppDB = database.New(o.Config.DB.Type, o.Config.DB.Sqlite.AppFilePath, database.WithMigrate(true))
		s.MonitorDB = database.New(o.Config.DB.Type, o.Config.DB.Sqlite.MonitorFilePath, database.WithMigrate(true))
		s.ScriptTaskDB = database.New(o.Config.DB.Type, o.Config.DB.Sqlite.ScriptTaskFilePath, database.WithMigrate(true))
	} else {
		Connect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			o.Config.DB.Mysql.Username,
			o.Config.DB.Mysql.Password,
			o.Config.DB.Mysql.Host,
			o.Config.DB.Mysql.Port,
			o.Config.DB.Mysql.DbName)
		s.AgentDB = database.New(o.Config.DB.Type, Connect, database.WithMigrate(true))
		s.AppDB = s.AgentDB
		s.MonitorDB = s.AgentDB
		s.ScriptTaskDB = s.AgentDB
	}

	s.Cron = cron.New(s.Config, s.AgentDB, s.AppDB, s.ScriptTaskDB, s.MonitorDB)

	// 初始化缓存
	s.Cache = o.initCache()

	return s, nil
}

// initCache 初始化缓存，默认使用内存方式
func (o *AppOptions) initCache() cache.Cache {
	cacheType := o.Config.Cache.Type
	if cacheType == "" {
		cacheType = "memory"
	}

	opts := []cache.Option{
		cache.WithMaxCost(o.Config.Cache.Memory.MaxCost),
		cache.WithBufferItems(o.Config.Cache.Memory.BufferItems),
		cache.WithMetrics(o.Config.Cache.Memory.Metrics),
	}

	connectionString := o.Config.Cache.GetConnectionString()

	c, err := cache.New(cacheType, connectionString, opts...)
	if err != nil {
		fmt.Printf("Cache initialization failed, falling back to memory cache: %v\n", err)
		// 降级到内存缓存
		c, _ = cache.New("memory", "", opts...)
	}

	return c
}

func (o *AppOptions) loadConfig(configFile string) {
	o.Config = config.New(configFile)
	if o.Config.Log.Path == "" {
		o.Config.Log.Path = "./log"
	}
	absPath, err := filepath.Abs(o.Config.Log.Path)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		absPath = o.Config.Log.Path
	}
	o.Config.Log.ErrorFilePath = absPath + "/" + o.Config.Log.ErrorFilename
	o.Config.Log.InfoFilePath = absPath + "/" + o.Config.Log.InfoFilename
}
