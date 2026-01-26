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

func (a *Application) List() response.Response {
	var applications []res.Application
	daoApps, err := a.Repository.List()
	if err != nil {
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
		return response.Error(model.ReturnErrCode(err))
	}

	// 如果应用正在运行，先停止服务
	if app.Status == "running" {
		stopRes := a.Stop(app.Name)
		if stopRes.Code != 200 {
			return stopRes
		}
	}

	// 删除数据库记录
	err = a.Repository.Delete(id)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

// DeleteByName 根据应用名称删除应用（用于回滚）
func (a *Application) DeleteByName(name string) response.Response {
	// 先获取应用信息
	apps, err := a.Repository.List()
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	var targetApp *model.Application
	for i := range apps {
		if apps[i].Name == name {
			targetApp = &apps[i]
			break
		}
	}

	if targetApp == nil {
		// 应用不存在，视为成功（幂等性）
		return response.Success("application not found, skip delete")
	}

	// 如果应用正在运行，先停止服务
	if targetApp.Status == "running" {
		stopRes := a.Stop(targetApp.Name)
		if stopRes.Code != 200 {
			zap.L().Warn("回滚时停止应用失败，继续删除",
				zap.Uint("id", targetApp.ID),
				zap.String("name", name),
			)
		}
	}

	// 删除数据库记录
	err = a.Repository.DeleteByName(name)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Add(request req.Application) response.Response {
	// 1. 检测 Docker 是否已安装
	if !checkDockerInstalled() {
		zap.L().Error("Docker 未安装")
		return response.Error(res.ErrDockerNotInstalled)
	}

	// 2. 检测 PATH 中是否有 docker-compose 命令
	if !checkDockerComposeAvailable() {
		zap.L().Error("docker-compose 命令未找到")
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
		zap.L().Error("创建 compose 目录失败", zap.Error(err))
		return response.Error(response.ErrSQL)
	}

	// 4. 创建 docker-compose 文件
	composeFileName := fmt.Sprintf("docker-compose-%s.yml", request.Name)
	composeFilePath := filepath.Join(composePath, composeFileName)

	// 如果文件已存在，先删除（支持重试）
	if _, err := os.Stat(composeFilePath); err == nil {
		zap.L().Info("docker-compose 文件已存在，先删除",
			zap.String("path", composeFilePath),
			zap.String("name", request.Name),
		)
		if err := os.Remove(composeFilePath); err != nil {
			zap.L().Error("删除已存在的 docker-compose 文件失败",
				zap.String("path", composeFilePath),
				zap.Error(err),
			)
			return response.Error(res.ErrDockerComposeCreate)
		}
	}

	if err := os.WriteFile(composeFilePath, []byte(request.Content), 0644); err != nil {
		zap.L().Error("创建 docker-compose 文件失败",
			zap.String("path", composeFilePath),
			zap.Error(err),
		)
		return response.Error(res.ErrDockerComposeCreate)
	}

	zap.L().Info("docker-compose 文件已创建",
		zap.String("path", composeFilePath),
		zap.String("name", request.Name),
	)

	// 5. 先保存到数据库，状态设置为 "starting"（启动中）
	modelReq := model.Application{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		Status:      "starting",
		Content:     request.Content,
		Version:     request.Version,
	}

	err := a.Repository.Add(&modelReq)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	confModel := model.Config{
		Key:   "server_id",
		Value: fmt.Sprint(request.ServerID),
	}
	err = a.ConfRepository.CreateOrUpdate(&confModel)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 6. 在协程中异步启动应用
	go func(appName, composePath, composeFileName string) {
		zap.L().Info("开始异步启动应用",
			zap.String("name", appName),
		)

		// 执行 docker-compose up -d 命令启动容器
		if err := runDockerComposeUp(composePath, composeFileName); err != nil {
			zap.L().Error("启动 docker-compose 失败",
				zap.String("path", composePath),
				zap.String("file", composeFileName),
				zap.Error(err),
			)
			// 启动失败，清理已创建的文件
			composeFilePath := filepath.Join(composePath, composeFileName)
			if removeErr := os.Remove(composeFilePath); removeErr != nil {
				zap.L().Error("启动失败后清理文件失败",
					zap.String("path", composeFilePath),
					zap.Error(removeErr),
				)
			}
			// 更新数据库状态为 "failed"
			apps, err := a.Repository.List()
			if err != nil {
				zap.L().Error("获取应用列表失败", zap.Error(err))
				return
			}
			for i := range apps {
				if apps[i].Name == appName {
					apps[i].Status = "failed"
					if updateErr := a.Repository.Update(&apps[i]); updateErr != nil {
						zap.L().Error("更新应用状态为 failed 失败", zap.Error(updateErr))
					}
					break
				}
			}
			return
		}

		zap.L().Info("docker-compose up 命令执行成功",
			zap.String("name", appName),
		)
		// 启动命令执行成功，但不立即更新数据库
		// 由 cron 定时任务 checkApplicationStatus 检测实际容器状态并更新
	}(request.Name, composePath, composeFileName)

	zap.L().Info("应用添加成功，正在后台启动",
		zap.String("name", request.Name),
	)

	return response.Success("应用添加成功，正在后台启动")
}

// getComposePath 获取 docker-compose 文件存储路径
func (a *Application) getComposePath() string {
	if a.Config.Common.ComposePath != "" {
		return a.Config.Common.ComposePath
	}
	return ""
}

// checkDockerInstalled 检测 Docker 是否已安装
func checkDockerInstalled() bool {
	_, err := execute.Command("docker", "--version")
	return err == nil
}

// checkDockerComposeAvailable 检测 docker-compose 命令是否可用
func checkDockerComposeAvailable() bool {
	// 检查 docker-compose 命令
	_, err := execute.Command("docker-compose", "--version")
	if err == nil {
		return true
	}

	// 检查 docker compose 插件
	_, err = execute.Command("docker", "compose", "version")
	if err == nil {
		return true
	}

	return false
}

// runDockerComposeUp 执行 docker-compose up -d 命令
func runDockerComposeUp(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "up", "-d"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "up", "-d"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose up -d",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
	)

	// 切换到工作目录执行命令
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	defer os.Chdir(currentDir)

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("切换到工作目录失败: %w", err)
	}

	// 执行命令并获取输出和错误
	stdout, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Error("docker-compose 命令执行失败",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose 启动成功",
		zap.String("output", stdout),
	)

	return nil
}

// checkCommandAvailable 检查命令是否在 PATH 中可用
func checkCommandAvailable(command string) bool {
	// 尝试执行命令的 --version 参数来检查命令是否可用
	_, err := execute.Command(command, "--version")
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
		return response.Error(model.ReturnErrCode(err))
	}

	return response.Success("success")
}

func (a *Application) Stop(name string) response.Response {
	// 获取应用信息
	app, err := a.Repository.GetByName(name)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != "running" {
		zap.L().Warn("应用未在运行中",
			zap.String("name", name),
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
	app.Status = "stopped"
	err = a.Repository.Update(&app)
	if err != nil {
		zap.L().Error("更新应用状态失败",
			zap.String("name", name),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	zap.L().Info("应用已停止",
		zap.String("name", name),
		zap.String("name", app.Name),
	)

	return response.Success("success")
}

// runDockerComposeStop 执行 docker-compose stop 命令
func runDockerComposeStop(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "stop"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "stop"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose stop",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
	)

	// 切换到工作目录执行命令
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	defer os.Chdir(currentDir)

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("切换到工作目录失败: %w", err)
	}

	// 执行命令并获取输出和错误
	stdout, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Error("docker-compose stop 命令执行失败",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose stop 执行成功",
		zap.String("output", stdout),
	)

	return nil
}

func (a *Application) Start(name string) response.Response {
	// 获取应用信息
	app, err := a.Repository.GetByName(name)
	if err != nil {
		return response.Error(model.ReturnErrCode(err))
	}

	// 检查应用状态
	if app.Status != "stopped" {
		zap.L().Warn("应用未处于停止状态",
			zap.String("name", name),
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
	app.Status = "starting"
	err = a.Repository.Update(&app)
	if err != nil {
		zap.L().Error("更新应用状态失败",
			zap.String("name", name),
			zap.Error(err),
		)
		return response.Error(model.ReturnErrCode(err))
	}

	// 在协程中异步启动应用
	go func(appName, composePath, composeFileName string) {
		zap.L().Info("开始异步启动应用",
			zap.String("name", appName),
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
				if apps[i].Name == appName {
					apps[i].Status = "failed"
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
		)
		// 启动命令执行成功，但不立即更新数据库
		// 由 cron 定时任务 checkApplicationStatus 检测实际容器状态并更新
	}(app.Name, composePath, composeFileName)

	zap.L().Info("启动请求已提交，正在后台处理",
		zap.String("name", name),
		zap.String("name", app.Name),
	)

	return response.Success("启动请求已提交，正在后台处理")
}

// runDockerComposeStart 执行 docker-compose start 命令
func runDockerComposeStart(workDir, composeFile string) error {
	var command string
	var args []string

	// 尝试使用 docker-compose 命令
	if checkCommandAvailable("docker-compose") {
		command = "docker-compose"
		args = []string{"-f", composeFile, "start"}
	} else if checkCommandAvailable("docker") {
		// 使用 docker compose 插件
		command = "docker"
		args = []string{"compose", "-f", composeFile, "start"}
	} else {
		return fmt.Errorf("docker-compose 命令不可用")
	}

	zap.L().Info("执行 docker-compose start",
		zap.String("work_dir", workDir),
		zap.String("compose_file", composeFile),
	)

	// 切换到工作目录执行命令
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	defer os.Chdir(currentDir)

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("切换到工作目录失败: %w", err)
	}

	// 执行命令并获取输出和错误
	stdout, stderr, err := execute.CommandError(command, args...)
	if err != nil {
		zap.L().Error("docker-compose start 命令执行失败",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info("docker-compose start 执行成功",
		zap.String("output", stdout),
	)

	return nil
}
