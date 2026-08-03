package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"squirrel-dev/internal/squ-agent/module/application/domain"
	"squirrel-dev/pkg/execute"
	"squirrel-dev/pkg/utils"
)

type composeProject struct {
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	ConfigFiles string `json:"ConfigFiles"`
}

func (j *Jobs) checkApplicationStatus() {
	ctx := context.Background()
	applications, err := j.applications.List(ctx)
	if err != nil {
		return
	}
	for _, app := range applications {
		status := j.containerStatus(app.DeployID)
		shouldUpdate := app.OldStatus != status || status == domain.StatusStarting || status == domain.StatusFailed
		if !shouldUpdate {
			continue
		}
		updated := app
		updated.Status = status
		updated.OldStatus = app.Status
		url := utils.GenAgentUrl(
			j.config.Apiserver.Http.Scheme,
			j.config.Apiserver.Http.Server,
			0,
			j.config.Apiserver.Http.BaseUri,
			uriAppReport,
		)
		_, _ = j.http.Post(url, applicationStatusReport{DeployID: updated.DeployID, Status: status}, nil)
		_ = j.applications.Update(ctx, &updated)
	}
}

func (j *Jobs) containerStatus(deployID uint64) string {
	status := j.checkComposeStatus("docker", "compose", deployID)
	if status != "" {
		return status
	}
	status = j.checkComposeStatus("docker-compose", "", deployID)
	if status != "" {
		return status
	}
	return domain.StatusFailed
}

func (j *Jobs) checkComposeStatus(command, prefix string, deployID uint64) string {
	var args []string
	if prefix == "" {
		args = []string{"ls", "--all", "--format", "json"}
	} else {
		args = []string{prefix, "ls", "--all", "--format", "json"}
	}
	output, _, err := execute.CommandError(command, args...)
	if err != nil {
		// Preserve the old behavior: this non-empty failure status means the
		// docker-compose fallback in containerStatus is effectively unreachable.
		return domain.StatusFailed
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &j.composeProjects); err != nil {
		return domain.StatusFailed
	}
	expected := fmt.Sprintf("%d/docker-compose.yml", deployID)
	for _, project := range j.composeProjects {
		if !strings.Contains(project.ConfigFiles, expected) {
			continue
		}
		status := strings.ToLower(project.Status)
		switch {
		case strings.HasPrefix(status, "running"):
			return domain.StatusRunning
		case strings.HasPrefix(status, "exited"):
			return domain.StatusStopped
		case strings.HasPrefix(status, "paused"):
			return domain.StatusPaused
		}
	}
	return domain.StatusUndeploy
}
