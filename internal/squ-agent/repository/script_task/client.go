package script_task

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func (c *Client) Add(task *model.ScriptExecutionTask) error {
	return c.DB.Create(task).Error
}

func (c *Client) Get(id uint) (model.ScriptExecutionTask, error) {
	var task model.ScriptExecutionTask
	err := c.DB.Where("id = ?", id).First(&task).Error
	return task, err
}

func (c *Client) GetRunningTask() (model.ScriptExecutionTask, error) {
	var task model.ScriptExecutionTask
	err := c.DB.Where("status = ?", "running").First(&task).Error
	return task, err
}

func (c *Client) Update(task *model.ScriptExecutionTask) error {
	return c.DB.Updates(task).Error
}

func (c *Client) GetUnreportedTasks() ([]model.ScriptExecutionTask, error) {
	var tasks []model.ScriptExecutionTask
	err := c.DB.Where("reported = ? AND status IN (?, ?)", false, "success", "failed").
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

func (c *Client) MarkAsReported(id uint) error {
	return c.DB.Model(&model.ScriptExecutionTask{}).
		Where("id = ?", id).
		Update("reported", true).Error
}
