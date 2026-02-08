package application

import (
	"fmt"
	"os"
	"path/filepath"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/application/res"
	"squirrel-dev/internal/squ-agent/model"

	"go.uber.org/zap"
)

func (a *Application) Start(deployID uint64) response.Response {
	// 获取应用信息
	app, err := a.Repository.GetByDeployID(deployID)
	if err != nil {
		zap.L().Error("Failed to get application by deploy id", zap.Uint64("deploy_id", deployID), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != model.AppStatusStopped {
		zap.L().Warn("Application is not in stopped state",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.String("status", app.Status),
		)
		return response.Error(res.ErrDockerComposeStart)
	}

	// 确定 docker-compose 文件路径
	composePath := a.getComposePath()
	if composePath == "" {
		composePath = "."
	}
	composeFileName := fmt.Sprintf("docker-compose-%s.yml", app.Name)
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 检查 docker-compose 文件是否存在
	if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
		zap.L().Error("docker-compose file not found",
			zap.String("path", composeFilePath),
		)
		return response.Error(res.ErrDockerComposeStart)
	}

	// 更新数据库状态为 "starting"（启动中）
	app.Status = model.AppStatusStarting
	err = a.Repository.Update(&app)
	if err != nil {
		zap.L().Error("Failed to update application status",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 在协程中异步启动应用
	go func(appName, composePath, composeFileName string, deployID uint64) {
		zap.L().Info("Starting application asynchronously",
			zap.String("name", appName),
			zap.Uint64("deploy_id", deployID),
		)

		// 执行 docker-compose start 命令
		if err := runDockerComposeStart(composePath, composeFileName); err != nil {
			zap.L().Error("Failed to start docker-compose",
				zap.String("path", composePath),
				zap.String("file", composeFileName),
				zap.Error(err),
			)
			// 启动失败，更新数据库状态
			apps, err := a.Repository.List()
			if err != nil {
				zap.L().Error("Failed to list applications", zap.Error(err))
				return
			}
			for i := range apps {
				if apps[i].DeployID == deployID {
					apps[i].Status = model.AppStatusFailed
					if updateErr := a.Repository.Update(&apps[i]); updateErr != nil {
						zap.L().Error("Failed to update application status to failed", zap.Error(updateErr))
					}
					break
				}
			}
			return
		}

		zap.L().Info("docker-compose start command executed successfully",
			zap.String("name", appName),
			zap.Uint64("deploy_id", deployID),
		)
		// 启动命令执行成功，但不立即更新数据库
		// 由 cron 定时任务 checkApplicationStatus 检测实际容器状态并更新
	}(app.Name, composePath, composeFileName, deployID)

	zap.L().Info("Start request submitted, processing in background",
		zap.String("name", app.Name),
		zap.Uint64("deploy_id", deployID),
	)

	return response.Success("Start request submitted, processing in background")
}

func (a *Application) Stop(deployID uint64) response.Response {
	// 获取应用信息
	app, err := a.Repository.GetByDeployID(deployID)
	if err != nil {
		zap.L().Error("Failed to get application by deploy id", zap.Uint64("deploy_id", deployID), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != model.AppStatusRunning {
		zap.L().Warn("Application is not running",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.String("status", app.Status),
		)
		return response.Error(res.ErrDockerComposeStop)
	}

	// 确定 docker-compose 文件路径
	composePath := a.getComposePath()
	if composePath == "" {
		composePath = "."
	}
	composeFileName := fmt.Sprintf("docker-compose-%s.yml", app.Name)
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 检查 docker-compose 文件是否存在
	if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
		zap.L().Error("docker-compose file not found",
			zap.String("path", composeFilePath),
		)
		return response.Error(res.ErrDockerComposeStop)
	}

	// 执行 docker-compose stop 命令
	if err := runDockerComposeStop(composePath, composeFileName); err != nil {
		zap.L().Error("Failed to stop docker-compose",
			zap.String("path", composePath),
			zap.String("file", composeFileName),
			zap.Error(err),
		)
		return response.Error(res.ErrDockerComposeStop)
	}

	// 更新数据库状态
	app.Status = model.AppStatusStopped
	err = a.Repository.Update(&app)
	if err != nil {
		zap.L().Error("Failed to update application status",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("Application stopped",
		zap.String("name", app.Name),
		zap.Uint64("deploy_id", deployID),
	)

	return response.Success("success")
}
