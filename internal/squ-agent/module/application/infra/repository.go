package infra

import (
	"context"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/application/domain"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.Application, error) {
	var models []applicationModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Application
	for _, model := range models {
		result = append(result, toDomain(model))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Application, error) {
	var model applicationModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Application{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) GetByDeployID(ctx context.Context, deployID uint64) (domain.Application, error) {
	var model applicationModel
	if err := r.db.WithContext(ctx).Where("deploy_id = ?", deployID).First(&model).Error; err != nil {
		return domain.Application{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&applicationModel{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Application) error {
	model := toModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Application) error {
	return r.db.WithContext(ctx).Updates(toModel(*value)).Error
}

func (r *Repository) Transaction(ctx context.Context, fn func(domain.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepository(tx))
	})
}

func toModel(value domain.Application) applicationModel {
	return applicationModel{
		ID: value.ID, Name: value.Name, Description: value.Description, Type: value.Type,
		OldStatus: value.OldStatus, Status: value.Status, Content: value.Content,
		Version: value.Version, DeployID: value.DeployID, Env: value.Env,
	}
}

func toDomain(value applicationModel) domain.Application {
	return domain.Application{
		ID: value.ID, Name: value.Name, Description: value.Description, Type: value.Type,
		OldStatus: value.OldStatus, Status: value.Status, Content: value.Content,
		Version: value.Version, DeployID: value.DeployID, Env: value.Env,
	}
}
