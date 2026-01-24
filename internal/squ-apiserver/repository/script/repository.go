package script

import (
	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

// ScriptRepository 定义脚本仓库的接口
type ScriptRepository interface {
	// List 获取所有脚本列表
	List() (scripts []model.Script, err error)

	// Get 根据 ID 获取单个脚本
	Get(id uint) (script model.Script, err error)

	// GetByName 根据名称获取脚本
	GetByName(name string) (script model.Script, err error)

	// Delete 根据 ID 删除脚本
	Delete(id uint) (err error)

	// Add 添加新脚本
	Add(req *model.Script) (err error)

	// Update 更新脚本
	Update(req *model.Script) (err error)

	// AddScriptResult 添加脚本执行结果
	AddScriptResult(result *model.ScriptResult) (err error)

	// GetScriptResults 获取脚本执行结果
	GetScriptResults(scriptID uint) (results []model.ScriptResult, err error)

	// GetLatestScriptResult 获取最新的脚本执行结果
	GetLatestScriptResult(scriptID uint, serverID uint) (result model.ScriptResult, err error)

	// UpdateScriptResult 更新脚本执行结果
	UpdateScriptResult(id uint, result *model.ScriptResult) (err error)
}

// New 创建新的脚本仓库实例
func New(db *gorm.DB) ScriptRepository {
	return &Client{
		DB: db,
	}
}
