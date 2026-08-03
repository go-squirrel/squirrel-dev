package infra

import (
	"context"
	"time"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/application/domain"
)

type model struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string
	Description string
	Type        string
	Content     string
	Version     string
}

func (model) TableName() string { return "applications" }

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.Application, error) {
	var models []model
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Application
	for _, value := range models {
		result = append(result, toDomain(value))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Application, error) {
	var value model
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&value).Error; err != nil {
		return domain.Application{}, err
	}
	return toDomain(value), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Application) error {
	record := toModel(*value)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return err
	}
	value.ID = record.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Application) error {
	return r.db.WithContext(ctx).Updates(toModel(*value)).Error
}

func toModel(v domain.Application) model {
	return model{ID: v.ID, Name: v.Name, Description: v.Description, Type: v.Type, Content: v.Content, Version: v.Version}
}

func toDomain(v model) domain.Application {
	return domain.Application{ID: v.ID, Name: v.Name, Description: v.Description, Type: v.Type, Content: v.Content, Version: v.Version}
}
