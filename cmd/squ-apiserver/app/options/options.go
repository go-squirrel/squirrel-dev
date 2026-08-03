package options

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/log"
	"squirrel-dev/internal/squ-apiserver/app"
	"squirrel-dev/internal/squ-apiserver/config"
)

// AppOptions 负责承接启动参数，并组装应用所需依赖。
type AppOptions struct {
	ConfFile string
	Config   *config.Config
}

func NewAppOptions() *AppOptions {
	return &AppOptions{}
}

// NewServer 兼容当前模板家族的命名习惯，返回装配后的 App。
func (o *AppOptions) NewServer() (*app.App, error) {
	instance := app.New()
	o.loadConfig(o.ConfFile)
	instance.Config = o.Config

	gin.SetMode(instance.Config.Server.Mode)
	instance.Gin = gin.New()

	instance.Log = log.NewClient(
		o.Config.Log.InfoFilePath,
		o.Config.Log.ErrorFilePath,
		o.Config.Log.Level,
		o.Config.Log.MaxSize,
		o.Config.Log.MaxBackups,
		o.Config.Log.MaxAge,
	)

	var connStr string
	if o.Config.DB.Type == "sqlite" {
		connStr = o.Config.DB.Sqlite.FilePath
	} else {
		connStr = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			o.Config.DB.Mysql.Username,
			o.Config.DB.Mysql.Password,
			o.Config.DB.Mysql.Host,
			o.Config.DB.Mysql.Port,
			o.Config.DB.Mysql.DbName,
		)
	}

	instance.DB = database.New(o.Config.DB.Type, connStr, database.WithMigrate(true))

	return instance, nil
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
