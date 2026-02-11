package application

import (
	"fmt"

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

	// 3. 检查 deployID 是否已存在，若存在则停止后重部署
	existingApp, err := a.Repository.GetByDeployID(request.DeployID)
	if err == nil {
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
		err = a.Repository.Transaction(func(appRepo appRepository.Repository) error {
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

	// 4. 新应用，正常部署流程
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Status:      model.AppStatusStopped, // 初始状态为 stopped
		Content:     request.Content,
		Version:     request.Version,
		DeployID:    request.DeployID,
	}

	err = a.Repository.Add(&modelReq)
	if err != nil {
		zap.L().Error("Failed to add application to database",
			zap.String("name", request.Name),
			zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 保存 server_id 配置
	confModel := model.Config{
		Key:   "server_id",
		Value: fmt.Sprint(request.ServerID),
	}
	if err := a.ConfRepository.CreateOrUpdate(&confModel); err != nil {
		zap.L().Error("Failed to create or update config", zap.Error(err))
		return response.Error(model.ReturnErrCode(err))
	}

	// 部署应用
	if _, _, err := a.deployApplication(&modelReq); err != nil {
		zap.L().Error("Failed to deploy application",
			zap.String("name", request.Name),
			zap.Error(err))
		return response.Error(res.ErrDockerComposeCreate)
	}

	zap.L().Info("Application added successfully, starting in background",
		zap.String("name", request.Name))

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
