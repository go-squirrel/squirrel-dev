package api

import (
	"squirrel-dev/internal/squ-apiserver/module/script/api/req"
	"squirrel-dev/internal/squ-apiserver/module/script/api/res"
	"squirrel-dev/internal/squ-apiserver/module/script/application"
	"squirrel-dev/internal/squ-apiserver/module/script/domain"
)

func toScriptRequest(value req.Script) application.ScriptRequest {
	return application.ScriptRequest{
		ID:      value.ID,
		Name:    value.Name,
		Content: value.Content,
	}
}

func toExecuteRequest(value req.ExecuteScript) application.ExecuteRequest {
	return application.ExecuteRequest{
		ScriptID: value.ScriptID,
		ServerID: value.ServerID,
	}
}

func toResultReport(value req.ScriptResultReport) application.ResultReport {
	return application.ResultReport{
		TaskID:       value.TaskID,
		ScriptID:     value.ScriptID,
		Output:       value.Output,
		Status:       value.Status,
		ErrorMessage: value.ErrorMessage,
	}
}

func toScriptResponse(value domain.Script) res.Script {
	return res.Script{
		ID:      value.ID,
		Name:    value.Name,
		Content: value.Content,
	}
}

func toScriptResultResponse(value domain.ScriptResult) res.ScriptResult {
	return res.ScriptResult{
		ID:           value.ID,
		TaskID:       value.TaskID,
		ScriptID:     value.ScriptID,
		ServerID:     value.ServerID,
		ServerIP:     value.ServerIP,
		AgentPort:    value.AgentPort,
		Output:       value.Output,
		Status:       value.Status,
		ErrorMessage: value.ErrorMessage,
		CreatedAt:    value.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
