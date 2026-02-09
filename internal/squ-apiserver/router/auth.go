package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/auth"
	"squirrel-dev/internal/squ-apiserver/handler/auth/res"
	authRepository "squirrel-dev/internal/squ-apiserver/repository/auth"
)

func Auth(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	res.RegisterCode()
	service := auth.Auth{
		Config:     conf,
		Repository: authRepository.New(db.GetDB()),
	}
	group.POST("/login", auth.LoginHandler(&service))
}
