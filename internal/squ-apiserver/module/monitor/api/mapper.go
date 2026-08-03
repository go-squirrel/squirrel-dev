package api

import (
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/module/monitor/domain"
)

func toResponse(value domain.Result) response.Response {
	return response.Response{
		Code:    response.CodeSuccess,
		Message: value.Message,
		Data:    value.Data,
	}
}
