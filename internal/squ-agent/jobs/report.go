package jobs

import (
	"context"
	"encoding/json"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/pkg/utils"
)

const (
	uriScriptResults = "/scripts/receive-result"
	uriAppReport     = "/deployment/report"
)

type applicationStatusReport struct {
	ApplicationID uint   `json:"application_id"`
	ServerID      uint   `json:"server_id"`
	Status        string `json:"status"`
	DeployID      uint64 `json:"deploy_id"`
}

type scriptResultReport struct {
	TaskID       uint   `json:"task_id"`
	ScriptsID    uint   `json:"script_id"`
	Output       string `json:"output"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func (j *Jobs) reportScriptResults() {
	ctx := context.Background()
	tasks, err := j.scriptTasks.GetUnreportedTasks(ctx)
	if err != nil {
		return
	}
	url := utils.GenAgentUrl(
		j.config.Apiserver.Http.Scheme,
		j.config.Apiserver.Http.Server,
		0,
		j.config.Apiserver.Http.BaseUri,
		uriScriptResults,
	)
	for _, task := range tasks {
		body, err := j.http.Post(url, scriptResultReport{
			TaskID: task.TaskID, ScriptsID: task.ScriptID, Output: task.Output,
			Status: task.Status, ErrorMessage: task.ErrorMsg,
		}, nil)
		if err != nil {
			continue
		}
		var result response.Response
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}
		if result.Code == 0 && task.Status == "success" {
			_ = j.scriptTasks.MarkAsReported(ctx, task.ID)
		}
	}
}
