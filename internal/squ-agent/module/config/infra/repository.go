package infra

import (
	"context"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-agent/module/config/domain"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]domain.Config, error) {
	var records []configModel
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}

	var configs []domain.Config
	for _, record := range records {
		configs = append(configs, toDomain(record))
	}
	return configs, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Config, error) {
	var record configModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		return domain.Config{}, err
	}
	return toDomain(record), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&configModel{}, id).Error
}

func (r *Repository) CreateOrUpdate(ctx context.Context, config *domain.Config) error {
	var record configModel
	err := r.db.WithContext(ctx).Where("key = ?", config.Key).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		record.Key = config.Key
		record.Value = config.Value
		if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
			return err
		}
		config.ID = record.ID
		return nil
	}

	record.Key = config.Key
	record.Value = config.Value
	config.ID = record.ID
	return r.db.WithContext(ctx).Updates(&record).Error
}

func (r *Repository) GetByKey(ctx context.Context, key string) (string, error) {
	var record configModel
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&record).Error; err != nil {
		return "", err
	}
	return record.Value, nil
}

func (r *Repository) Transaction(ctx context.Context, fn func(domain.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepository(tx))
	})
}

func toDomain(record configModel) domain.Config {
	return domain.Config{
		ID:    record.ID,
		Key:   record.Key,
		Value: record.Value,
	}
}
