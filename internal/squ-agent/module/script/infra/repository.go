package infra

import (
	"context"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/script/domain"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Add(ctx context.Context, task *domain.Task) error {
	model := toModel(*task)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	task.ID = model.ID
	return nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Task, error) {
	var model taskModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Task{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) GetRunningTask(ctx context.Context) (domain.Task, error) {
	var model taskModel
	if err := r.db.WithContext(ctx).Where("status = ?", "running").First(&model).Error; err != nil {
		return domain.Task{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) Update(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Updates(toModel(*task)).Error
}

func (r *Repository) GetUnreportedTasks(ctx context.Context) ([]domain.Task, error) {
	var models []taskModel
	err := r.db.WithContext(ctx).
		Where("reported = ? AND status IN (?, ?)", false, "success", "failed").
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	var result []domain.Task
	for _, model := range models {
		result = append(result, toDomain(model))
	}
	return result, nil
}

func (r *Repository) MarkAsReported(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&taskModel{}).Where("id = ?", id).Update("reported", true).Error
}

func toModel(value domain.Task) taskModel {
	return taskModel{
		ID: value.ID, ScriptID: value.ScriptID, TaskID: value.TaskID, Name: value.Name,
		Content: value.Content, Status: value.Status, Output: value.Output, ErrorMsg: value.ErrorMsg,
		ExecutedAt: value.ExecutedAt, Reported: value.Reported,
	}
}

func toDomain(value taskModel) domain.Task {
	return domain.Task{
		ID: value.ID, ScriptID: value.ScriptID, TaskID: value.TaskID, Name: value.Name,
		Content: value.Content, Status: value.Status, Output: value.Output, ErrorMsg: value.ErrorMsg,
		ExecutedAt: value.ExecutedAt, Reported: value.Reported,
	}
}
