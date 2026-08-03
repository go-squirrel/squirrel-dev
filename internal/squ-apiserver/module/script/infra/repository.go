package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/script/domain"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.Script, error) {
	var models []scriptModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Script
	for _, model := range models {
		result = append(result, toDomainScript(model))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Script, error) {
	var model scriptModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Script{}, err
	}
	return toDomainScript(model), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	var model scriptModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Script) error {
	var existing scriptModel
	result := r.db.WithContext(ctx).Where("name = ?", value.Name).First(&existing)
	if result.Error == nil {
		return errors.New("script with this name already exists")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	model := toScriptModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Script) error {
	var existing scriptModel
	if err := r.db.WithContext(ctx).First(&existing, value.ID).Error; err != nil {
		return err
	}
	if value.Name != existing.Name {
		var conflict scriptModel
		result := r.db.WithContext(ctx).Where("name = ? AND id != ?", value.Name, value.ID).First(&conflict)
		if result.Error == nil {
			return errors.New("script with this name already exists")
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
	}
	return r.db.WithContext(ctx).Updates(toScriptModel(*value)).Error
}

func (r *Repository) AddResult(ctx context.Context, value *domain.ScriptResult) error {
	model := toResultModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	value.CreatedAt = model.CreatedAt
	return nil
}

func (r *Repository) ListResults(ctx context.Context, scriptID uint) ([]domain.ScriptResult, error) {
	var models []resultModel
	if err := r.db.WithContext(ctx).Where("script_id = ?", scriptID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.ScriptResult
	for _, model := range models {
		result = append(result, toDomainResult(model))
	}
	return result, nil
}

func (r *Repository) UpdateResultByTaskID(ctx context.Context, taskID uint64, value *domain.ScriptResult) error {
	return r.db.WithContext(ctx).Model(&resultModel{}).Where("task_id = ?", taskID).Updates(toResultModel(*value)).Error
}

func toDomainScript(value scriptModel) domain.Script {
	return domain.Script{ID: value.ID, Name: value.Name, Content: value.Content}
}

func toScriptModel(value domain.Script) scriptModel {
	return scriptModel{ID: value.ID, Name: value.Name, Content: value.Content}
}

func toDomainResult(value resultModel) domain.ScriptResult {
	return domain.ScriptResult{
		ID: value.ID, CreatedAt: value.CreatedAt, TaskID: value.TaskID, ScriptID: value.ScriptID,
		ServerID: value.ServerID, ServerIP: value.ServerIP, AgentPort: value.AgentPort,
		Output: value.Output, Status: value.Status, ErrorMessage: value.ErrorMessage,
	}
}

func toResultModel(value domain.ScriptResult) resultModel {
	return resultModel{
		ID: value.ID, TaskID: value.TaskID, ScriptID: value.ScriptID, ServerID: value.ServerID,
		ServerIP: value.ServerIP, AgentPort: value.AgentPort, Output: value.Output,
		Status: value.Status, ErrorMessage: value.ErrorMessage,
	}
}
