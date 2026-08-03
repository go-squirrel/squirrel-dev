package infra

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/server/domain"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"
)

type AgentClient struct {
	config *config.Config
	client *httpclient.Client
}

func NewAgentClient(conf *config.Config, client *httpclient.Client) *AgentClient {
	return &AgentClient{config: conf, client: client}
}

func (c *AgentClient) GetInfo(_ context.Context, ip string, port int) (string, map[string]any) {
	url := utils.GenAgentUrl(
		c.config.Agent.Http.Scheme,
		ip,
		port,
		c.config.Agent.Http.BaseUrl,
		"server/info",
	)
	start := time.Now()
	body, err := c.client.Get(url, nil)
	if err != nil {
		zap.L().Error("failed to get agent information",
			zap.String("url", url),
			zap.String("method", "GET"),
			zap.String("ip_address", ip),
			zap.Int("agent_port", port),
			zap.Duration("cost", time.Since(start)),
			zap.Error(err),
		)
		return domain.StatusOffline, nil
	}
	var result response.Response
	if err := json.Unmarshal(body, &result); err != nil {
		zap.L().Error("failed to parse agent response",
			zap.String("url", url),
			zap.String("method", "GET"),
			zap.String("ip_address", ip),
			zap.Int("agent_port", port),
			zap.Duration("cost", time.Since(start)),
			zap.Error(err),
		)
		return domain.StatusOffline, nil
	}
	if result.Code != 0 {
		zap.L().Error("agent failed to get server information",
			zap.String("url", url),
			zap.String("method", "GET"),
			zap.String("ip_address", ip),
			zap.Int("agent_port", port),
			zap.Int("code", result.Code),
			zap.String("message", result.Message),
			zap.Duration("cost", time.Since(start)),
		)
		return domain.StatusOffline, nil
	}
	if result.Data == nil {
		return domain.StatusOnline, nil
	}
	info, _ := result.Data.(map[string]any)
	return domain.StatusOnline, info
}
