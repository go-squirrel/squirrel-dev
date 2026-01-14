package router

import (
	"github.com/gin-gonic/gin"

	"squirrel-dev/internal/pkg/database"
	"squirrel-dev/internal/squ-apiserver/config"
	appStoreHandler "squirrel-dev/internal/squ-apiserver/handler/app_store"

	appStoreModel "squirrel-dev/internal/squ-apiserver/model/app_store"
)

func AppStore(group *gin.RouterGroup, conf *config.Config, db database.DB) {
	service := appStoreHandler.AppStore{
		Config:      conf,
		ModelClient: appStoreModel.New(db.GetDB()),
	}
	group.GET("/app-store", appStoreHandler.ListHandler(&service))
	group.GET("/app-store/:id", appStoreHandler.GetHandler(&service))
	group.DELETE("/app-store/:id", appStoreHandler.DeleteHandler(&service))
	group.POST("/app-store", appStoreHandler.AddHandler(&service))
	group.POST("/app-store/:id", appStoreHandler.UpdateHandler(&service))
}
