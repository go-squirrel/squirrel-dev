package api

import "squirrel-dev/internal/pkg/response"

func toResponse(value any) response.Response {
	return response.Success(value)
}
