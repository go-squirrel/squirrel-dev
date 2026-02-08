package server

import (
	"embed"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/pkg/middleware/cors"
	"squirrel-dev/internal/pkg/middleware/log"
	spaMiddleware "squirrel-dev/internal/pkg/middleware/static"
	"squirrel-dev/internal/squ-apiserver/config"
)

//go:embed all:dist
var staticData embed.FS

type Server struct {
	Config *config.Config
	Gin    *gin.Engine
	// 导入日志
	Log *log.Client
	DB  database.DB
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
	staticFunc, err := static.EmbedFolder(staticData, "dist")
	if err != nil {
		zap.L().Error("failed to embed static folder",
			zap.String("folder", "dist"),
			zap.Error(err),
		)
	}
	s.Gin.Use(static.Serve("/", staticFunc))

	// SPA 路由回退中间件：处理前端路由
	s.Gin.Use(spaMiddleware.Default(staticData, "dist"))

	// s.Gin.Use(log.GinLogger(s.Log.Logger),
	// 	log.GinRecovery(s.Log.Logger, true),
	// 	c)
	s.Gin.Use(c)

	s.migrate()

	s.SetupRouter()

	err = s.Gin.Run(s.Config.Server.Bind + ":" + s.Config.Server.Port)
	if err != nil {
		zap.L().Error("failed to start server",
			zap.String("address", s.Config.Server.Bind+":"+s.Config.Server.Port),
			zap.Error(err),
		)
	}
}
