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
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != model.AppStatusStopped {
		zap.L().Warn("应用未处于停止状态",
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
		zap.L().Error("docker-compose 文件不存在",
			zap.String("path", composeFilePath),
		)
		return response.Error(res.ErrDockerComposeStart)
	}

	// 更新数据库状态为 "starting"（启动中）
	app.Status = model.AppStatusStarting
	err = a.Repository.Update(&app)
	if err != nil {
		zap.L().Error("更新应用状态失败",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 在协程中异步启动应用
	go func(appName, composePath, composeFileName string, deployID uint64) {
		zap.L().Info("开始异步启动应用",
			zap.String("name", appName),
			zap.Uint64("deploy_id", deployID),
		)

		// 执行 docker-compose start 命令
		if err := runDockerComposeStart(composePath, composeFileName); err != nil {
			zap.L().Error("启动 docker-compose 失败",
				zap.String("path", composePath),
				zap.String("file", composeFileName),
				zap.Error(err),
			)
			// 启动失败，更新数据库状态
			apps, err := a.Repository.List()
			if err != nil {
				zap.L().Error("获取应用列表失败", zap.Error(err))
				return
			}
			for i := range apps {
				if apps[i].DeployID == deployID {
					apps[i].Status = model.AppStatusFailed
					if updateErr := a.Repository.Update(&apps[i]); updateErr != nil {
						zap.L().Error("更新应用状态为 failed 失败", zap.Error(updateErr))
					}
					break
				}
			}
			return
		}

		zap.L().Info("docker-compose start 命令执行成功",
			zap.String("name", appName),
			zap.Uint64("deploy_id", deployID),
		)
		// 启动命令执行成功，但不立即更新数据库
		// 由 cron 定时任务 checkApplicationStatus 检测实际容器状态并更新
	}(app.Name, composePath, composeFileName, deployID)

	zap.L().Info("启动请求已提交，正在后台处理",
		zap.String("name", app.Name),
		zap.Uint64("deploy_id", deployID),
	)

	return response.Success("启动请求已提交，正在后台处理")
}

func (a *Application) Stop(deployID uint64) response.Response {
	// 获取应用信息
	app, err := a.Repository.GetByDeployID(deployID)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != model.AppStatusRunning {
		zap.L().Warn("应用未在运行中",
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
		zap.L().Error("docker-compose 文件不存在",
			zap.String("path", composeFilePath),
		)
		return response.Error(res.ErrDockerComposeStop)
	}

	// 执行 docker-compose stop 命令
	if err := runDockerComposeStop(composePath, composeFileName); err != nil {
		zap.L().Error("停止 docker-compose 失败",
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
		zap.L().Error("更新应用状态失败",
			zap.String("name", app.Name),
			zap.Uint64("deploy_id", deployID),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("应用已停止",
		zap.String("name", app.Name),
		zap.Uint64("deploy_id", deployID),
	)

	return response.Success("success")
}
