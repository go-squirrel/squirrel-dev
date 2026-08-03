package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/module/script/domain"
	serverDomain "squirrel-dev/internal/squ-apiserver/module/server/domain"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"
)

type ServerReader struct {
	repository serverDomain.Repository
}

func NewServerReader(repository serverDomain.Repository) *ServerReader {
	return &ServerReader{repository: repository}
}

func (r *ServerReader) Get(ctx context.Context, id uint) (domain.Server, error) {
	server, err := r.repository.Get(ctx, id)
	if err != nil {
		return domain.Server{}, err
	}
	return domain.Server{ID: server.ID, IPAddress: server.IPAddress, AgentPort: server.AgentPort}, nil
}

type AgentClient struct {
	config *config.Config
	http   *httpclient.Client
}

func NewAgentClient(conf *config.Config) *AgentClient {
	return &AgentClient{config: conf, http: httpclient.NewClient(30 * time.Second)}
}

func (c *AgentClient) Post(_ context.Context, server domain.Server, path string, request any) error {
	url := utils.GenAgentUrl(
		c.config.Agent.Http.Scheme,
		server.IPAddress,
		server.AgentPort,
		c.config.Agent.Http.BaseUrl,
		path,
	)
	start := time.Now()
	logger := zap.L().With(
		zap.String("url", url),
		zap.String("method", "POST"),
		zap.Uint("server_id", server.ID),
		zap.String("agent_path", path),
	)
	body, err := c.http.Post(url, request, nil)
	if err != nil {
		logger.Error("agent request failed", zap.Duration("cost", time.Since(start)), zap.Error(err))
		return fmt.Errorf("agent request failed: %w", err)
	}
	var result response.Response
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("failed to parse agent response", zap.Duration("cost", time.Since(start)), zap.Error(err))
		return fmt.Errorf("parse agent response failed: %w", err)
	}
	if result.Code != 0 {
		logger.Error("agent returned error",
			zap.Int("code", result.Code),
			zap.String("message", result.Message),
			zap.Duration("cost", time.Since(start)),
		)
		return fmt.Errorf("agent error: code=%d, message=%s", result.Code, result.Message)
	}
	return nil
}

type IDGenerator struct{}

func (IDGenerator) Generate() (uint64, error) { return utils.IDGenerate() }
