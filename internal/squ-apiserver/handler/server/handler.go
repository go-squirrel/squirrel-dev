package server

import (
	"log"
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func ListHandler(service *Server) func(c *gin.Context) {

	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}

func GetHandler(service *Server) func(c *gin.Context) {

	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Get(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func DeleteHandler(service *Server) func(c *gin.Context) {

	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Delete(idUint)
		c.JSON(http.StatusOK, res)
	}
}

func AddHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Server{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Add(request)
		c.JSON(http.StatusOK, res)
	}
}

func UpdateHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request := req.Server{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		request.ID = idUint
		res := service.Update(request)
		c.JSON(http.StatusOK, res)
	}
}

func RegistryHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Register{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.Registry(request)
		c.JSON(http.StatusOK, res)
	}
}

func TerminalHandler(service *Server) func(c *gin.Context) {

	return func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket升级失败: %v", err)
			return
		}

		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.S().Warn(err)
			c.JSON(http.StatusBadRequest, response.Error(response.ErrCodeParameter))
			return
		}
		res := service.GetTerminal(idUint, conn)

		conn.WriteJSON(res)
	}
}
