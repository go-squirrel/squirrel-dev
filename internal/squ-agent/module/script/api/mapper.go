package api

import (
	"squirrel-dev/internal/squ-agent/module/script/api/req"
	"squirrel-dev/internal/squ-agent/module/script/application"
)

func toExecuteRequest(value req.ExecuteScript) application.Request {
	return application.Request{
		ID:      value.ID,
		Name:    value.Name,
		Content: value.Content,
		TaskID:  value.TaskID,
	}
}
