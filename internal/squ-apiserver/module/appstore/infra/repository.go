package infra

import (
	"context"
	"time"

	"gorm.io/gorm"

	"squirrel-dev/internal/squ-apiserver/module/appstore/domain"
)

type appModel struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"column:name;type:varchar(100);not null;unique;comment:模板名称"`
	Description string         `gorm:"column:description;type:text;comment:模板描述"`
	Type        string         `gorm:"column:type;type:varchar(50);not null;comment:类型(compose/k8s_manifest/helm_chart)"`
	Category    string         `gorm:"column:category;type:varchar(50);comment:分类(web/database/middleware/devops)"`
	Icon        *string        `gorm:"column:icon;type:varchar(255);comment:图标URL"`
	Version     string         `gorm:"column:version;type:varchar(50);not null;comment:版本"`
	Content     string         `gorm:"column:content;type:text;not null;comment:模板内容"`
	Tags        string         `gorm:"column:tags;type:varchar(255);comment:标签(逗号分隔)"`
	Author      string         `gorm:"column:author;type:varchar(100);comment:作者"`
	RepoURL     *string        `gorm:"column:repo_url;type:varchar(255);comment:仓库地址"`
	HomepageURL *string        `gorm:"column:homepage_url;type:varchar(255);comment:主页地址"`
	IsOfficial  bool           `gorm:"column:is_official;type:tinyint(1);default:false;comment:是否官方模板"`
	Downloads   int            `gorm:"column:downloads;type:int(10);default:0;comment:下载次数"`
	Status      string         `gorm:"column:status;type:varchar(20);default:'active';comment:状态(active/deprecated)"`
}

func (appModel) TableName() string { return "app_stores" }

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]domain.App, error) {
	var models []appModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []domain.App
	for _, model := range models {
		result = append(result, toDomain(model))
	}
	return result, nil
}

func (r *Repository) Get(ctx context.Context, id uint) (domain.App, error) {
	var model appModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return domain.App{}, err
	}
	return toDomain(model), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&appModel{}, id).Error
}

func (r *Repository) Add(ctx context.Context, value *domain.App) error {
	model := toModel(*value)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	value.ID = model.ID
	return nil
}

func (r *Repository) Update(ctx context.Context, value *domain.App) error {
	return r.db.WithContext(ctx).Updates(toModel(*value)).Error
}

func toModel(v domain.App) appModel {
	return appModel{ID: v.ID, Name: v.Name, Description: v.Description, Type: v.Type, Category: v.Category, Icon: v.Icon, Version: v.Version, Content: v.Content, Tags: v.Tags, Author: v.Author, RepoURL: v.RepoURL, HomepageURL: v.HomepageURL, IsOfficial: v.IsOfficial, Downloads: v.Downloads, Status: v.Status}
}

func toDomain(v appModel) domain.App {
	return domain.App{ID: v.ID, Name: v.Name, Description: v.Description, Type: v.Type, Category: v.Category, Icon: v.Icon, Version: v.Version, Content: v.Content, Tags: v.Tags, Author: v.Author, RepoURL: v.RepoURL, HomepageURL: v.HomepageURL, IsOfficial: v.IsOfficial, Downloads: v.Downloads, Status: v.Status}
}
