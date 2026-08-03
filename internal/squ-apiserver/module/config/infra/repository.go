package infra

import (
	"context"
	"time"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/config/domain"
)

type configModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Key       string
	Value     string
}

func (configModel) TableName() string { return "configs" }

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.Config, error) {
	var models []configModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Config
	for _, model := range models {
		result = append(result, domain.Config{ID: model.ID, Key: model.Key, Value: model.Value})
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Config, error) {
	var model configModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Config{}, err
	}
	return domain.Config{ID: model.ID, Key: model.Key, Value: model.Value}, nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&configModel{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Config) error {
	model := configModel{Key: value.Key, Value: value.Value}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Config) error {
	return r.db.WithContext(ctx).Updates(&configModel{ID: value.ID, Key: value.Key, Value: value.Value}).Error
}
