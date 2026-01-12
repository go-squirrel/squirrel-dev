package auth

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"

	"squirrel-dev/internal/squ-apiserver/handler/auth/req"
	"squirrel-dev/internal/squ-apiserver/handler/auth/res"
	authModel "squirrel-dev/internal/squ-apiserver/model/auth"
	"time"

	"squirrel-dev/pkg/jwt"

	"go.uber.org/zap"
)

type Auth struct {
	Config      *config.Config
	ModelClient authModel.Repository
}

func (a *Auth) Login(request req.Request) response.Response {

	ok := a.ModelClient.Get(request.Username, request.Password)
	if !ok {
		return response.Error(response.ErrUserOrPassword)
	}

	j := jwt.New(a.Config.Auth.Jwt.SigningKey)
	expireDuration := time.Duration(a.Config.Auth.Jwt.Expired) * time.Minute

	token, err := j.GenToken(request.Username, expireDuration)
	if err != nil {
		zap.S().Error(err)
	}

	return response.Success(res.TokenRes{
		Token: token,
	})
}
