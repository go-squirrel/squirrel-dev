package application

import (
	"fmt"
	"os"
	"path/filepath"

	"squirrel-dev/internal/pkg/response"
	"squirrel-dev/internal/squ-agent/config"
	"squirrel-dev/internal/squ-agent/handler/application/req"
	"squirrel-dev/internal/squ-agent/handler/application/res"
	"squirrel-dev/internal/squ-agent/model"
	"squirrel-dev/pkg/execute"

	appRepository "squirrel-dev/internal/squ-agent/repository/application"
	confRepository "squirrel-dev/internal/squ-agent/repository/config"

	"go.uber.org/zap"
)

type Application struct {
	Config         *config.Config
	Repository     appRepository.Repository
	ConfRepository confRepository.Repository
}

func New(config *config.Config, repo appRepository.Repository, confRepo confRepository.Repository) *Application {
	return &Application{
		Config:         config,
		Repository:     repo,
		ConfRepository: confRepo,
	}
}

func (a *Application) List() response.Response {
	var applications []res.Application
	daoApps, err := a.Repository.List()
	if err != nil {
		zap.L().Error("Failed to list applications", zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}
	for _, daoA := range daoApps {
		applications = append(applications, res.Application{
			ID:          daoA.ID,
			Name:        daoA.Name,
			Description: daoA.Description,
			Type:        daoA.Type,
			Status:      daoA.Status,
			Content:     daoA.Content,
			Version:     daoA.Version,
		})
	}
	return response.Success(applications)
}

func (a *Application) Get(id uint) response.Response {
	var appRes res.Application
	daoA, err := a.Repository.Get(id)
	if err != nil {
		zap.L().Error("Failed to get application", zap.Uint("id", id), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}
	appRes = res.Application{
		ID:          daoA.ID,
		Name:        daoA.Name,
		Description: daoA.Description,
		Type:        daoA.Type,
		Status:      daoA.Status,
		Content:     daoA.Content,
		Version:     daoA.Version,
	}

	return response.Success(appRes)
}

func (a *Application) Delete(id uint) response.Response {
	// 先获取应用信息
	app, err := a.Repository.Get(id)
	if err != nil {
		zap.L().Error("Failed to get application", zap.Uint("id", id), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 如果应用正在运行，先停止服务
	if app.Status == model.AppStatusRunning {
		stopRes := a.Stop(app.DeployID)
		if stopRes.Code != 200 {
			return stopRes
		}
	}

	// 删除数据库记录
	err = a.Repository.Delete(id)
	if err != nil {
		zap.L().Error("Failed to delete application", zap.Uint("id", id), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

// DeleteByDeployID 根据deployID删除应用
func (a *Application) DeleteByDeployID(deployID uint64) response.Response {
	// 先获取应用信息
	app, err := a.Repository.GetByDeployID(deployID)
	if err != nil {
		if err.Error() == "record not found" {
			// 应用不存在，视为成功（幂等性）
			return response.Success("application not found, skip delete")
		}
		zap.L().Error("Failed to get application by deploy id", zap.Uint64("deploy_id", deployID), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 如果应用正在运行，先停止服务
	if app.Status == model.AppStatusRunning {
		stopRes := a.Stop(deployID)
		if stopRes.Code != 200 {
			zap.L().Warn("Failed to stop application during deletion, continuing with deletion",
				zap.Uint("id", app.ID),
				zap.String("name", app.Name),
				zap.Uint64("deploy_id", deployID),
			)
		}
	}

	// 删除数据库记录
	err = a.Repository.Delete(app.ID)
	if err != nil {
		zap.L().Error("Failed to delete application", zap.Uint("id", app.ID), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Add(request req.Application) response.Response {
	// 1. 检测 Docker 是否已安装
	if !checkDockerInstalled() {
		zap.L().Error("Docker is not installed")
		return response.Error(res.ErrDockerNotInstalled)
	}

	// 2. 检测 PATH 中是否有 docker-compose 命令
	if !checkDockerComposeAvailable() {
		zap.L().Error("docker-compose command not found")
		return response.Error(res.ErrDockerComposeNotFound)
	}

	// 3. 确定 docker-compose 文件存储路径
	composePath := a.getComposePath()
	if composePath == "" {
		// 如果配置文件中没有设置，使用当前目录
		composePath = "."
	}

	// 确保目录存在
	if err := os.MkdirAll(composePath, 0755); err != nil {
		zap.L().Error("Failed to create compose directory", zap.String("path", composePath), zap.Error(err))
		return response.Error(response.ErrSQL)
	}

	// 4. 创建 docker-compose 文件
	composeFileName := fmt.Sprintf("docker-compose-%s.yml", request.Name)
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 如果文件已存在，先删除（支持重试）
	if _, err := os.Stat(composeFilePath); err == nil {
		zap.L().Info("docker-compose file already exists, deleting it", zap.String("path", composeFilePath), zap.String("name", request.Name))
		if err := os.Remove(composeFilePath); err != nil {
			zap.L().Error("Failed to delete existing docker-compose file", zap.String("path", composeFilePath), zap.Error(err))
			return response.Error(res.ErrDockerComposeCreate)
		}
	}

	if err := os.WriteFile(composeFilePath, []byte(request.Content), 0644); err != nil {
		zap.L().Error("Failed to create docker-compose file", zap.String("path", composeFilePath), zap.Error(err))
		return response.Error(res.ErrDockerComposeCreate)
	}

	zap.L().Info("docker-compose file created", zap.String("path", composeFilePath), zap.String("name", request.Name))

	// 5. 先保存到数据库，状态设置为 "starting"（启动中）
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Status:      model.AppStatusStarting,
		Content:     request.Content,
		Version:     request.Version,
		DeployID:    request.DeployID,
	}

	err := a.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("Failed to add application to database", zap.String("name", request.Name), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	confModel := model.Config{
		Key:   "server_id",
		Value: fmt.Sprint(request.ServerID),
	}
	err = a.ConfRepository.CreateOrUpdate(&confModel)
	if err != nil {
		zap.L().Error("Failed to create or update config", zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 6. 在协程中异步启动应用
	go func(appName, composePath, composeFileName string) {
		zap.L().Info("Starting application asynchronously", zap.String("name", appName))

		// 执行 docker-compose up -d 命令启动容器
		if err := runDockerComposeUp(composePath, composeFileName); err != nil {
			zap.L().Error("Failed to start docker-compose", zap.String("path", composePath), zap.String("file", composeFileName), zap.Error(err))
			// 更新数据库状态为 "failed"
			apps, err := a.Repository.List()
			if err != nil {
				zap.L().Error("Failed to list applications", zap.Error(err))
				return
			}
			for i := range apps {
				if apps[i].Name == appName {
					apps[i].Status = model.AppStatusFailed
					if updateErr := a.Repository.Update(&apps[i]); updateErr != nil {
						zap.L().Error("Failed to update application status to failed", zap.Error(updateErr))
					}
					break
				}
			}
			return
		}

		zap.L().Info("docker-compose up command executed successfully",
			zap.String("name", appName),
		)
		// 启动命令执行成功，但不立即更新数据库
		// 由 cron 定时任务 checkApplicationStatus 检测实际容器状态并更新
	}(request.Name, composePath, composeFileName)

	zap.L().Info("Application added successfully, starting in background", zap.String("name", request.Name))

	return response.Success("Application added successfully, starting in background")
}

// checkDockerInstalled 检测 Docker 是否已安装
func checkDockerInstalled() bool {
	_, err := execute.Command("docker", "--version")
	return err == nil
}

func (a *Application) Update(request req.Application) response.Response {
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Status:      request.Status,
		Content:     request.Content,
		Version:     request.Version,
	}
	modelReq.ID = request.ID
	err := a.Repository.Update(&modelReq)

	if err != nil {
		zap.L().Error("Failed to update application", zap.Uint("id", request.ID), zap.String("name", request.Name), zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}
