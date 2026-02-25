package server

import (
	"time"

	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/cache"
	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/cors"
	"squirrel-dev/internal/pkg/middleware/log"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/cron"

	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Gin    *gin.Engine
	// 导入日志
	Log          *log.Client
	AgentDB      database.DB
	AppDB        database.DB
	MonitorDB    database.DB
	ScriptTaskDB database.DB
	Cache        cache.Cache
	Cron         *cron.Cron
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	c := cors.New(cors.Config{
		AllowOrigins:     s.Config.Server.Origins,
		AllowMethods:     s.Config.Server.Methods,
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})

	s.Gin.Use(log.GinLogger(s.Log.Logger),
		log.GinRecovery(s.Log.Logger, true),
		c)

	s.migrate()

	s.SetupRouter()

	s.Cron.Start()

	err := s.Gin.Run(s.Config.Server.Bind + ":" + s.Config.Server.Port)
	if err != nil {
		zap.L().Error("Failed to start server", zap.Error(err))
	}
}
