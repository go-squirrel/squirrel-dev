package script_task

import (
	"squirrel-dev/internal/squ-agent/model"

	"gorm.io/gorm"
)

type Repository interface {
	Add(task *model.ScriptExecutionTask) (err error)
	Get(id uint) (task model.ScriptExecutionTask, err error)
	GetRunningTask() (task model.ScriptExecutionTask, err error)
	Update(task *model.ScriptExecutionTask) (err error)
	GetUnreportedTasks() (tasks []model.ScriptExecutionTask, err error)
	MarkAsReported(id uint) (err error)
}

func New(db *gorm.DB) Repository {
	return &Client{
		DB: db,
	}
}
