package infra

import (
	"context"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/server/domain"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.Server, error) {
	var models []serverModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Server
	for _, model := range models {
		result = append(result, toDomain(model))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Server, error) {
	var model serverModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Server{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&serverModel{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Server) error {
	model := toModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Server) error {
	return r.db.WithContext(ctx).Updates(toModel(*value)).Error
}

func (r *Repository) GetByUUID(ctx context.Context, uuid string) (domain.Server, error) {
	var model serverModel
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&model).Error; err != nil {
		return domain.Server{}, err
	}
	return toDomain(model), nil
}

func toModel(value domain.Server) serverModel {
	return serverModel{
		ID: value.ID, UUID: value.UUID, Hostname: value.Hostname, IPAddress: value.IPAddress,
		AgentPort: value.AgentPort, SSHUsername: value.SSHUsername, SSHPassword: value.SSHPassword,
		SSHPrivateKey: value.SSHPrivateKey, SSHPassphrase: value.SSHPassphrase, SSHPort: value.SSHPort,
		AuthType: value.AuthType, ServerAlias: value.ServerAlias, Status: value.Status,
	}
}

func toDomain(value serverModel) domain.Server {
	return domain.Server{
		ID: value.ID, UUID: value.UUID, Hostname: value.Hostname, IPAddress: value.IPAddress,
		AgentPort: value.AgentPort, SSHUsername: value.SSHUsername, SSHPassword: value.SSHPassword,
		SSHPrivateKey: value.SSHPrivateKey, SSHPassphrase: value.SSHPassphrase, SSHPort: value.SSHPort,
		AuthType: value.AuthType, ServerAlias: value.ServerAlias, Status: value.Status,
	}
}
