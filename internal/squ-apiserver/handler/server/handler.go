package server

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/server/req"
	"squirrel-dev/internal/squ-apiserver/handler/server/res"
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
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
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
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
			return
		}
		resp := service.Delete(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func AddHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Server{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
			return
		}
		resp := service.Add(request)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
			return
		}
		request := req.Server{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
			return
		}
		request.ID = idUint
		resp := service.Update(request)
		c.JSON(http.StatusOK, resp)
	}
}

func RegistryHandler(service *Server) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Register{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidParameter))
			return
		}
		resp := service.Registry(request)
		c.JSON(http.StatusOK, resp)
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
			zap.L().Error("failed to upgrade websocket connection",
				zap.Error(err),
			)
			return
		}

		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			conn.WriteJSON(response.Error(res.ErrInvalidParameter))
			conn.Close()
			return
		}
		resp := service.GetTerminal(idUint, conn)

		conn.WriteJSON(resp)
	}
}
