package script

import (
	"net/http"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/script/req"
	"squirrel-dev/internal/squ-apiserver/handler/script/res"
	"squirrel-dev/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := service.List()
		c.JSON(http.StatusOK, res)
	}
}

func GetHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.Get(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.Delete(idUint)
		c.JSON(http.StatusOK, resp)
	}
}

func AddHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.Script{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.Add(request)
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		request := req.Script{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		request.ID = idUint
		resp := service.Update(request)
		c.JSON(http.StatusOK, resp)
	}
}

func ExecuteHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.ExecuteScript{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.Execute(request)
		c.JSON(http.StatusOK, resp)
	}
}

func ReceiveResultHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := req.ScriptResultReport{}
		err := c.ShouldBindJSON(&request)
		if err != nil {
			zap.L().Warn("failed to bind request JSON",
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.ReceiveScriptResult(request)
		c.JSON(http.StatusOK, resp)
	}
}

func GetResultsHandler(service *Script) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := utils.StringToUint(id)
		if err != nil {
			zap.L().Warn("failed to convert id to uint",
				zap.String("id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusOK, response.Error(res.ErrInvalidScriptContent))
			return
		}
		resp := service.GetResults(idUint)
		c.JSON(http.StatusOK, resp)
	}
}
