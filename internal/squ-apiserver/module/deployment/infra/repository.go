package infra

import (
	"context"
	"time"

	"gorm.io/gorm"

	applicationDomain "squirrel-dev/internal/squ-apiserver/module/application/domain"
	"squirrel-dev/internal/squ-apiserver/module/deployment/domain"
	serverDomain "squirrel-dev/internal/squ-apiserver/module/server/domain"
)

type deploymentModel struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ServerID      uint
	ApplicationID uint
	Status        string
	DeployID      uint64
	Content       string
	Env           []map[string]string `gorm:"type:json;serializer:json"`
}

func (deploymentModel) TableName() string { return "deployments" }

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context, serverID uint) ([]domain.Deployment, error) {
	var models []deploymentModel
	query := r.db.WithContext(ctx)
	if serverID > 0 {
		query = query.Where("server_id = ?", serverID)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.Deployment
	for _, model := range models {
		result = append(result, toDomain(model))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.Deployment, error) {
	var model deploymentModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.Deployment{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) GetByDeployID(ctx context.Context, deployID uint64) (domain.Deployment, error) {
	var model deploymentModel
	if err := r.db.WithContext(ctx).Where("deploy_id = ?", deployID).First(&model).Error; err != nil {
		return domain.Deployment{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&deploymentModel{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.Deployment) error {
	model := toModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID, value.CreatedAt = model.ID, model.CreatedAt
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.Deployment) error {
	return r.db.WithContext(ctx).Updates(toModel(*value)).Error
}

func (r *Repository) UpdateStatus(ctx context.Context, deployID uint64, status string) error {
	return r.db.WithContext(ctx).Model(&deploymentModel{}).Where("deploy_id = ?", deployID).Update("status", status).Error
}

func toModel(v domain.Deployment) deploymentModel {
	return deploymentModel{ID: v.ID, CreatedAt: v.CreatedAt, ServerID: v.ServerID,
		ApplicationID: v.ApplicationID, Status: v.Status, DeployID: v.DeployID, Content: v.Content, Env: v.Env}
}

func toDomain(v deploymentModel) domain.Deployment {
	return domain.Deployment{ID: v.ID, CreatedAt: v.CreatedAt, ServerID: v.ServerID,
		ApplicationID: v.ApplicationID, Status: v.Status, DeployID: v.DeployID, Content: v.Content, Env: v.Env}
}

type ApplicationReader struct{ repository applicationDomain.Repository }

func NewApplicationReader(repository applicationDomain.Repository) *ApplicationReader {
	return &ApplicationReader{repository: repository}
}

func (r *ApplicationReader) Get(ctx context.Context, id uint) (domain.Application, error) {
	value, err := r.repository.Get(ctx, id)
	if err != nil {
		return domain.Application{}, err
	}
	return domain.Application{ID: value.ID, Name: value.Name, Description: value.Description,
		Type: value.Type, Content: value.Content, Version: value.Version}, nil
}

type ServerReader struct{ repository serverDomain.Repository }

func NewServerReader(repository serverDomain.Repository) *ServerReader {
	return &ServerReader{repository: repository}
}

func (r *ServerReader) Get(ctx context.Context, id uint) (domain.Server, error) {
	value, err := r.repository.Get(ctx, id)
	if err != nil {
		return domain.Server{}, err
	}
	return domain.Server{ID: value.ID, IPAddress: value.IPAddress, AgentPort: value.AgentPort}, nil
}
