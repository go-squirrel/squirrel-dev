package application

import (
	"fmt"
	"os"
	"path/filepath"
	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/handler/application/req"
	"squirrel-dev/internal/squ-agent/handler/application/res"
	"squirrel-dev/internal/squ-agent/model"
	appRepository "squirrel-dev/internal/squ-agent/repository/application"

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
	composePath := a.getComposePathOrDefault()
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
			a.updateApplicationStatusToFailed(deployID)
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
	composePath := a.getComposePathOrDefault()
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

func (a *Application) redeploy(existingApp model.Application, request req.Application) response.Response {
	zap.L().Info("Application with deployID already exists, will redeploy",
		zap.Uint64("deploy_id", request.DeployID),
		zap.String("existing_name", existingApp.Name),
		zap.String("new_name", request.Name))

	// 如果应用正在运行，先停止
	if existingApp.Status == model.AppStatusRunning {
		zap.L().Info("Stopping existing application before redeploy", zap.Uint64("deploy_id", request.DeployID))
		if stopErr := runDockerComposeStop(a.getComposePathOrDefault(),
			fmt.Sprintf("docker-compose-%s.yml", existingApp.Name)); stopErr != nil {
			zap.L().Warn("Failed to stop existing application, continuing with redeploy", zap.Uint64("deploy_id", request.DeployID), zap.Error(stopErr))
		}
	}

	// 使用事务更新应用信息和保存配置
	err := a.Repository.Transaction(func(appRepo appRepository.Repository) error {
		// 更新现有应用的信息
		existingApp.Name = request.Name
		existingApp.Description = request.Description
		existingApp.Type = request.Type
		existingApp.Content = request.Content
		existingApp.Version = request.Version
		existingApp.Status = model.AppStatusStarting // 更新状态为 starting
		if err := appRepo.Update(&existingApp); err != nil {
			return fmt.Errorf("failed to update application: %w", err)
		}

		// 保存 server_id 配置
		confModel := model.Config{
			Key:   "server_id",
			Value: fmt.Sprint(request.ServerID),
		}
		return a.ConfRepository.CreateOrUpdate(&confModel)
	})
	if err != nil {
		zap.L().Error("Failed to update application and config in transaction", zap.Uint64("deploy_id", request.DeployID), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 准备 compose 文件
	composePath, composeFileName, err := a.prepareComposeFiles(&existingApp)
	if err != nil {
		zap.L().Error("Failed to prepare compose files", zap.Uint64("deploy_id", request.DeployID), zap.Error(err))
		return response.Error(res.ErrDockerComposeCreate)
	}

	// 异步启动应用
	a.startDockerComposeAsync(request.Name, composePath, composeFileName, request.DeployID)

	zap.L().Info("Application redeployed successfully, starting in background", zap.String("name", request.Name), zap.Uint64("deploy_id", request.DeployID))

	return response.Success("Application redeployed successfully, starting in background")
}
