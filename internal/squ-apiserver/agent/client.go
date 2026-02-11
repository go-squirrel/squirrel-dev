package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/config"
	"squirrel-dev/internal/squ-apiserver/model"
	"squirrel-dev/pkg/httpclient"
	"squirrel-dev/pkg/utils"

	"go.uber.org/zap"
)

// Client 统一的 Agent 客户端
type Client struct {
	config     *config.Config
	httpClient *httpclient.Client
}

// NewClient 创建新的 Agent 客户端
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config:     cfg,
		httpClient: httpclient.NewClient(30 * time.Second),
	}
}

// CallResult 封装 Agent 调用结果
type CallResult struct {
	Resp   response.Response
	URL    string
	Err    error
	ErrCode int // 业务错误码，用于上层返回
}

// Call 统一处理 Agent 调用，包含完整的错误处理和日志记录
// method: GET 或 POST
// path: Agent 的 API 路径
// body: 请求体（POST 时使用，GET 时可传 nil）
// logFields: 额外的日志字段
func (c *Client) Call(ctx context.Context, server model.Server, method, path string, body any, logFields ...zap.Field) CallResult {
	url := utils.GenAgentUrl(
		c.config.Agent.Http.Scheme,
		server.IpAddress,
		server.AgentPort,
		c.config.Agent.Http.BaseUrl,
		path,
	)

	// 基础日志字段
	fields := append([]zap.Field{
		zap.String("url", url),
		zap.String("method", method),
		zap.String("server_ip", server.IpAddress),
		zap.Int("agent_port", server.AgentPort),
	}, logFields...)

	var respBody []byte
	var err error

	switch method {
	case "GET":
		respBody, err = c.httpClient.Get(url, nil)
	case "POST":
		respBody, err = c.httpClient.Post(url, body, nil)
	default:
		return CallResult{
			URL: url,
			Err: fmt.Errorf("unsupported method: %s", method),
		}
	}

	if err != nil {
		zap.L().Error("agent request failed", append(fields, zap.Error(err))...)
		return CallResult{
			URL: url,
			Err: fmt.Errorf("agent request failed: %w", err),
		}
	}

	var agentResp response.Response
	if err := json.Unmarshal(respBody, &agentResp); err != nil {
		zap.L().Error("failed to parse agent response", append(fields, zap.Error(err))...)
		return CallResult{
			URL: url,
			Err: fmt.Errorf("parse agent response failed: %w", err),
		}
	}

	if agentResp.Code != 0 {
		zap.L().Error("agent returned error",
			append(fields,
				zap.Int("code", agentResp.Code),
				zap.String("message", agentResp.Message),
			)...)
		return CallResult{
			Resp: agentResp,
			URL:  url,
			Err:  fmt.Errorf("agent error: code=%d, message=%s", agentResp.Code, agentResp.Message),
		}
	}

	zap.L().Debug("agent call success", fields...)
	return CallResult{
		Resp: agentResp,
		URL:  url,
	}
}

// CallSimple 简化版调用，直接返回 response.Response 和 error
func (c *Client) CallSimple(ctx context.Context, server model.Server, method, path string, body any, logFields ...zap.Field) (response.Response, error) {
	result := c.Call(ctx, server, method, path, body, logFields...)
	return result.Resp, result.Err
}

// Get 发送 GET 请求
func (c *Client) Get(ctx context.Context, server model.Server, path string, logFields ...zap.Field) CallResult {
	return c.Call(ctx, server, "GET", path, nil, logFields...)
}

// Post 发送 POST 请求
func (c *Client) Post(ctx context.Context, server model.Server, path string, body any, logFields ...zap.Field) CallResult {
	return c.Call(ctx, server, "POST", path, body, logFields...)
}
