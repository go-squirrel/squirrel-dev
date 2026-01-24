package script

import (
	"errors"

	"squirrel-dev/internal/squ-apiserver/model"

	"gorm.io/gorm"
)

// Client 脚本仓库客户端实现
type Client struct {
	DB *gorm.DB
}

// List 获取所有脚本列表
func (c *Client) List() (scripts []model.Script, err error) {
	err = c.DB.Find(&scripts).Error
	return scripts, err
}

// Get 根据 ID 获取单个脚本
func (c *Client) Get(id uint) (script model.Script, err error) {
	err = c.DB.Where("id = ?", id).First(&script).Error
	if err != nil {
		return script, err
	}
	return script, nil
}

// GetByName 根据名称获取脚本
func (c *Client) GetByName(name string) (script model.Script, err error) {
	err = c.DB.Where("name = ?", name).First(&script).Error
	if err != nil {
		return script, err
	}
	return script, nil
}

// Delete 根据 ID 删除脚本
func (c *Client) Delete(id uint) (err error) {
	// 先检查脚本是否存在
	var script model.Script
	result := c.DB.First(&script, id)
	if result.Error != nil {
		return result.Error
	}
	return c.DB.Delete(&script).Error
}

// Add 添加新脚本
func (c *Client) Add(req *model.Script) (err error) {
	// 检查同名脚本是否已存在
	var existingScript model.Script
	result := c.DB.Where("name = ?", req.Name).First(&existingScript)
	if result.Error == nil {
		return errors.New("script with this name already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return c.DB.Create(req).Error
}

// Update 更新脚本
func (c *Client) Update(req *model.Script) (err error) {
	// 检查脚本是否存在
	var existingScript model.Script
	result := c.DB.First(&existingScript, req.ID)
	if result.Error != nil {
		return result.Error
	}

	// 如果更新了名称，检查新名称是否已被其他脚本使用
	if req.Name != existingScript.Name {
		var nameConflict model.Script
		result := c.DB.Where("name = ? AND id != ?", req.Name, req.ID).First(&nameConflict)
		if result.Error == nil {
			return errors.New("script with this name already exists")
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
	}

	return c.DB.Updates(req).Error
}

// AddScriptResult 添加脚本执行结果
func (c *Client) AddScriptResult(result *model.ScriptResult) (err error) {
	return c.DB.Create(result).Error
}

// GetScriptResults 获取脚本执行结果
func (c *Client) GetScriptResults(scriptID uint) ([]model.ScriptResult, error) {
	var results []model.ScriptResult
	err := c.DB.Where("script_id = ?", scriptID).Order("created_at DESC").Find(&results).Error
	return results, err
}

// GetLatestScriptResult 获取最新的脚本执行结果
func (c *Client) GetLatestScriptResult(scriptID, serverID uint) (result model.ScriptResult, err error) {
	err = c.DB.Where("script_id = ? AND server_id = ?", scriptID, serverID).
		Order("created_at DESC").
		First(&result).Error
	return result, err
}

// UpdateScriptResultByTaskID 根据 TaskID 更新脚本执行结果
func (c *Client) UpdateScriptResultByTaskID(taskID uint64, result *model.ScriptResult) (err error) {
	return c.DB.Model(&model.ScriptResult{}).Where("task_id = ?", taskID).Updates(result).Error
}
