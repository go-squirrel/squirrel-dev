package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	authModel "squirrel-dev/internal/squ-apiserver/model/auth"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/handler/auth"
)

func Auth(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	service := auth.Auth{
		Config: conf,
		ModelClient: authModel.New(db.GetDB()),
	}
	group.POST("/login", auth.LoginHandler(&service))
}
