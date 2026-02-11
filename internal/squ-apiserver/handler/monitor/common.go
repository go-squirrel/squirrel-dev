package monitor

import (
	"context"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-apiserver/handler/monitor/res"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (m *Monitor) callAgent(serverID uint, path string, description string) response.Response {
	server, err := m.ServerRepo.Get(serverID)
	if err != nil {
		zap.L().Error("failed to get server for monitoring",
			zap.Uint("server_id", serverID),
			zap.String("description", description),
			zap.Error(err),
		)
		return response.Error(returnMonitorErrCode(err))
	}

	result := m.AgentClient.Get(context.Background(), server, path,
		zap.Uint("server_id", serverID),
		zap.String("description", description),
	)
	if result.Err != nil {
		return response.Error(res.ErrMonitorFailed)
	}

	return result.Resp
}

// returnMonitorErrCode 根据错误类型返回精确的监控错误码
func returnMonitorErrCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return res.ErrServerNotFound
	}
	return res.ErrMonitorFailed
}
